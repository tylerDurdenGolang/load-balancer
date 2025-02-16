package healthcheck

import (
	"net"
	"time"

	"github.com/tylerDurdenGolang/load-balancer/internal/balancer"
)

type HealthChecker struct {
	interval time.Duration
	bal      balancer.IBalancer
}

func NewHealthChecker(interval time.Duration, b balancer.IBalancer) *HealthChecker {
	return &HealthChecker{
		interval: interval,
		bal:      b,
	}
}

func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		// Проверяем все бэкенды
		for _, backend := range hc.bal.GetAllBackends() {
			if hc.check(backend) {
				// «живой» бэкенд
				hc.bal.MarkBackendUp(backend)
			} else {
				// «упавший» бэкенд
				hc.bal.MarkBackendDown(backend)
			}
		}
	}
}

func (hc *HealthChecker) check(backend string) bool {
	// Пример простой проверки TCP-порта
	conn, err := net.DialTimeout("tcp", backend, 1*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
