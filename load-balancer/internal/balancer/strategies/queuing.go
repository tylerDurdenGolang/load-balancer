package strategies

// import (
// 	"math"
// 	"sync"

// 	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
// 	"github.com/tylerDurdenGolang/load-balancer/internal/errs"
// )

// type QueuingBalancer struct {
// 	backends []*domain.Backend
// 	lambda   float64 // Средняя интенсивность входящих запросов
// 	mutex    sync.Mutex
// }

// func NewQueuingBalancer(backends []*domain.Backend, lambda float64) *QueuingBalancer {
// 	return &QueuingBalancer{
// 		backends: backends,
// 		lambda:   lambda,
// 	}
// }

// func factorial(n int) float64 {
// 	result := 1.0
// 	for i := 2; i <= n; i++ {
// 		result *= float64(i)
// 	}
// 	return result
// }

// func (q *QueuingBalancer) erlangC(mu float64, c int) float64 {
// 	rho := q.lambda / (mu * float64(c))
// 	if rho >= 1 {
// 		return 1.0
// 	}

// 	term := math.Pow(q.lambda/mu, float64(c)) / (factorial(c) * (1 - rho))
// 	sum := 0.0
// 	for k := 0; k < c; k++ {
// 		sum += math.Pow(q.lambda/mu, float64(k)) / factorial(k)
// 	}

// 	return term / (sum + term)
// }

// func (q *QueuingBalancer) GetBackend() (*domain.Backend, error) {
// 	q.mutex.Lock()
// 	defer q.mutex.Unlock()

// 	minLoad := math.MaxFloat64
// 	var selected *domain.Backend

// 	for _, b := range q.backends {
// 		// b.Lock()
// 		if b.IsAlive() {
// 			mu := 1.0 / b.Metrics.Latency.Seconds()
// 			load := q.erlangC(mu, 1) // c=1 для одного сервера
// 			if load < minLoad {
// 				minLoad = load
// 				selected = b
// 			}
// 		}
// 		// b.Unlock()
// 	}

// 	if selected == nil {
// 		return nil, errs.ErrNoHealthyBackends
// 	}
// 	return selected, nil
// }
