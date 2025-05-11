// loadtest.go
//
// Простой генератор нагрузки для MathCruncher.
// ~10 000 RPS за счёт 1 000 горутин * 10 запросов в секунду каждая.
// Запуск:  go run loadtest.go -addr http://localhost:5001 -duration 30s
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	addr     = flag.String("addr", "http://localhost:8080", "base address of worker")
	duration = flag.Duration("duration", 30*time.Second, "test duration")
	rps      = flag.Int("rps", 10000, "target requests per second")
)

const payload = `{"expression":"sin(x)","lower":0,"upper":6.28,"samples":1000}`

func main() {
	flag.Parse()

	// Настройка HTTP-клиента (keep-alive, ограничение коннектов).
	tr := &http.Transport{
		MaxIdleConnsPerHost: 2000,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	// Рассчитываем воркеров и интервал.
	workers := 1000
	perWorker := *rps / workers                        // запрос/сек на воркер
	interval := time.Second / time.Duration(perWorker) // пауза между запросами

	var okCnt, errCnt uint64

	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	log.Printf("starting load: %d workers × %d req/sec each (interval %v)",
		workers, perWorker, interval)

	// Стартуем workers
	for i := 0; i < workers; i++ {
		go func() {
			tick := time.NewTicker(interval)
			defer tick.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tick.C:
					req, _ := http.NewRequest("POST", *addr+"/integrate",
						bytes.NewBufferString(payload))
					req.Header.Set("Content-Type", "application/json")
					resp, err := client.Do(req)
					if err != nil {
						atomic.AddUint64(&errCnt, 1)
						continue
					}
					_ = resp.Body.Close()
					if resp.StatusCode == 200 {
						atomic.AddUint64(&okCnt, 1)
					} else {
						atomic.AddUint64(&errCnt, 1)
					}
				}
			}
		}()
	}

	<-ctx.Done()
	total := atomic.LoadUint64(&okCnt) + atomic.LoadUint64(&errCnt)
	fmt.Println("------ result ------")
	fmt.Printf("duration:     %v\n", *duration)
	fmt.Printf("total sent:   %d\n", total)
	fmt.Printf("success 200:  %d\n", okCnt)
	fmt.Printf("errors:       %d\n", errCnt)
	fmt.Printf("achieved rps: %.0f\n",
		float64(total)/duration.Seconds())
}
