package main

import (
	"fmt"
	"log"

	"github.com/tylerDurdenGolang/load-balancer/internal/balancer"
	"github.com/tylerDurdenGolang/load-balancer/internal/config"
	"github.com/tylerDurdenGolang/load-balancer/internal/healthcheck"
	"github.com/tylerDurdenGolang/load-balancer/internal/server"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Инициализация балансировщика (например, Round Robin)
	lbStrategy, err := balancer.NewStrategy(cfg)
	if err != nil {
		log.Fatalf("failed to init strategy: %v", err)
	}
	// 3. Запуск health-check
	hc := healthcheck.NewHealthChecker(cfg.HealthCheckInterval, lbStrategy)
	go hc.Start()

	// 4. Старт HTTP-сервера
	svr := server.NewServer(cfg.ListenAddr, lbStrategy)
	fmt.Printf("Load Balancer is running on %s...\n", cfg.ListenAddr)

	if err := svr.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
