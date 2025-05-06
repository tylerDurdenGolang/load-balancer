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

// integrate performs Monte‑Carlo estimation of ∫_lower^upper f(x) dx
// returning mean (estimate) and 95 % confidence interval.
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
			n += remainder // add leftover samples
		}
		go func(n int) {
			localSum, localSq := 0.0, 0.0
			params := map[string]interface{}{}
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			for i := 0; i < n; i++ {
				x := lower + rng.Float64()*(upper-lower)
				params["x"] = x
				vRaw, err := expr.Evaluate(params)
				if err != nil {
					ch <- partial{}
					return
				}
				v := vRaw.(float64)
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
	variance := (totalSq/float64(samples) - mean*mean)
	ci95 = 1.96 * math.Sqrt(variance/float64(samples))
	return
}

func integrateHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() { reqLatency.Observe(time.Since(start).Seconds()) }()
	reqCounter.Inc()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Samples <= 0 {
		req.Samples = 100000
	}

	expr, err := govaluate.NewEvaluableExpression("(" + req.Expression + ")")
	if err != nil {
		http.Error(w, "invalid expression: "+err.Error(), http.StatusBadRequest)
		return
	}

	mean, ci95, err := integrate(expr, req.Lower, req.Upper, req.Samples)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response{Mean: mean, CI95: ci95, Samples: req.Samples}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/integrate", integrateHandler)
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("MathCruncher listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
