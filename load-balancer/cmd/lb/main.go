package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tylerDurdenGolang/load-balancer/internal/balancer"
	"github.com/tylerDurdenGolang/load-balancer/internal/config"
	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
	"github.com/tylerDurdenGolang/load-balancer/internal/healthcheck"
	"github.com/tylerDurdenGolang/load-balancer/internal/server"
)

func main() {
	/* ---------- 1. Загрузка конфигурации ---------- */

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	/* ---------- 2. Список backend‑ов на старте ---------- */

	backends := resolveBackends("template-api", cfg.WorkerPort)
	if len(backends) == 0 {
		log.Println("warning: no backends found on startup")
	}

	/* ---------- 3. Стратегия балансировки ---------- */

	lb, err := balancer.NewStrategy(cfg.Algorithm, backends)
	if err != nil {
		log.Fatalf("balancer init error: %v", err)
	}

	/* ---------- 4. Health‑checker ---------- */

	hc := healthcheck.NewHealthChecker(cfg.HealthCheckInterval, lb)
	go hc.Start()

	/* ---------- 5. Динамический DNS‑refresh ---------- */

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			newBack := resolveBackends("template-api", cfg.WorkerPort)
			lb.ReplaceBackends(newBack) // метод реализован в стратегии
			log.Printf("backend list refreshed: %d pods", len(newBack))
		}
	}()

	/* ---------- 6. HTTP‑сервер ---------- */

	srv := server.NewServer(cfg.ListenAddr, lb)
	log.Printf("Load Balancer listening on %s", cfg.ListenAddr)

	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

/*
resolveBackends выполняет DNS lookup headless‑service’а

	и собирает []*domain.Backend c учётом порта воркера.
*/
func resolveBackends(service string, port int) []*domain.Backend {
	ips, err := net.LookupHost(service)
	if err != nil {
		log.Printf("DNS lookup failed for %s: %v", service, err)
		return nil
	}

	backends := make([]*domain.Backend, 0, len(ips))
	for _, ip := range ips {
		addr := fmt.Sprintf("%s:%d", ip, port)
		log.Printf("discovered backend: %s", addr)
		backends = append(backends, domain.NewBackend(addr))
	}
	return backends
}
