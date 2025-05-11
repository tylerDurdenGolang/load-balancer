package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/* --------- типы запроса / ответа --------- */

type request struct {
	Expression string  `json:"expression"`
	Lower      float64 `json:"lower"`
	Upper      float64 `json:"upper"`
	Samples    int     `json:"samples"`
}

type response struct {
	Mean    float64 `json:"mean"`
	CI95    float64 `json:"ci95"`
	Samples int     `json:"samples"`
}

/* --------- Prometheus‑метрики --------- */

var (
	reqCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mc_requests_total",
		Help: "Total integrate requests",
	})
	reqLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "mc_req_latency_seconds",
		Help:    "Latency of integrate endpoint",
		Buckets: prometheus.DefBuckets,
	})
)

func init() {
	prometheus.MustRegister(reqCounter, reqLatency)
	rand.Seed(time.Now().UnixNano())
}

/* --------- карта поддерживаемых функций --------- */

var funcs = map[string]govaluate.ExpressionFunction{
	"sin":  func(a ...interface{}) (interface{}, error) { return math.Sin(a[0].(float64)), nil },
	"cos":  func(a ...interface{}) (interface{}, error) { return math.Cos(a[0].(float64)), nil },
	"tan":  func(a ...interface{}) (interface{}, error) { return math.Tan(a[0].(float64)), nil },
	"exp":  func(a ...interface{}) (interface{}, error) { return math.Exp(a[0].(float64)), nil },
	"sqrt": func(a ...interface{}) (interface{}, error) { return math.Sqrt(a[0].(float64)), nil },
	"log":  func(a ...interface{}) (interface{}, error) { return math.Log(a[0].(float64)), nil },
	"pow": func(a ...interface{}) (interface{}, error) {
		return math.Pow(a[0].(float64), a[1].(float64)), nil
	},
}

/* --------- Monte‑Carlo интегрирование --------- */

func integrate(expr *govaluate.EvaluableExpression, lower, upper float64, samples int) (mean, ci95 float64, err error) {
	workers := runtime.NumCPU()
	if samples < workers {
		workers = 1
	}
	chunk := samples / workers
	remainder := samples % workers

	type partial struct{ sum, sumSq float64 }
	ch := make(chan partial, workers)

	for w := 0; w < workers; w++ {
		n := chunk
		if w == 0 {
			n += remainder
		}
		go func(n int) {
			localSum, localSq := 0.0, 0.0
			params := map[string]interface{}{}
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			for i := 0; i < n; i++ {
				x := lower + rng.Float64()*(upper-lower)
				params["x"] = x
				val, _ := expr.Evaluate(params) // ошибок быть не должно
				v := val.(float64)
				localSum += v
				localSq += v * v
			}
			ch <- partial{localSum, localSq}
		}(n)
	}

	totalSum, totalSq := 0.0, 0.0
	for i := 0; i < workers; i++ {
		p := <-ch
		totalSum += p.sum
		totalSq += p.sumSq
	}

	mean = totalSum / float64(samples)
	variance := totalSq/float64(samples) - mean*mean
	ci95 = 1.96 * math.Sqrt(variance/float64(samples))
	return
}

/* --------- HTTP‑обработчики --------- */

func integrateHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	reqCounter.Inc()

	log.Printf("/integrate called from %s", r.RemoteAddr)

	defer func() {
		lat := time.Since(start).Seconds()
		reqLatency.Observe(lat)
		log.Printf("/integrate completed in %.3fs", lat)
	}()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("invalid JSON: %v", err)
		return
	}
	if req.Samples <= 0 {
		req.Samples = 100_000
	}

	log.Printf("expression='%s' lower=%.4f upper=%.4f samples=%d",
		req.Expression, req.Lower, req.Upper, req.Samples)

	expr, err := govaluate.NewEvaluableExpressionWithFunctions("("+req.Expression+")", funcs)
	if err != nil {
		http.Error(w, "invalid expression: "+err.Error(), http.StatusBadRequest)
		log.Printf("invalid expression: %v", err)
		return
	}

	mean, ci95, _ := integrate(expr, req.Lower, req.Upper, req.Samples)
	resp := response{Mean: mean, CI95: ci95, Samples: req.Samples}

	log.Printf("result: mean=%.6f ci95=%.6f", mean, ci95)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// Если нужно более сложное условие (например, проверка подключения к БД),
	// добавьте его здесь. Сейчас просто говорим «я жив и готов».
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/integrate", integrateHandler)
	mux.HandleFunc("/ready", readyHandler)
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("MathCruncher worker listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
