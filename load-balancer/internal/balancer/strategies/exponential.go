package strategies

// import (
// 	"math"
// 	"sync"

// 	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
// 	"github.com/tylerDurdenGolang/load-balancer/internal/errs"
// )

// type ExponentialBalancer struct {
// 	backends    []*domain.Backend
// 	alpha       float64
// 	predictions map[string]float64
// 	mutex       sync.RWMutex
// }

// func NewExponentialBalancer(backends []*domain.Backend, alpha float64) *ExponentialBalancer {
// 	return &ExponentialBalancer{
// 		backends:    backends,
// 		alpha:       alpha,
// 		predictions: make(map[string]float64),
// 	}
// }

// func (e *ExponentialBalancer) predict(b *domain.Backend) float64 {
// 	current := 0.4*b.Metrics.CPU + 0.3*b.Metrics.MemUsage + 0.3*float64(b.Metrics.Latency.Milliseconds())
// 	prev, exists := e.predictions[b.Host()]
// 	if !exists {
// 		return current
// 	}
// 	return e.alpha*current + (1-e.alpha)*prev
// }

// func (e *ExponentialBalancer) GetBackend() (*domain.Backend, error) {
// 	e.mutex.Lock()
// 	defer e.mutex.Unlock()

// 	minPrediction := math.MaxFloat64
// 	var selected *domain.Backend

// 	for _, b := range e.backends {
// 		// b.Lock()
// 		if b.IsAlive() {
// 			pred := e.predict(b)
// 			e.predictions[b.Host] = pred
// 			if pred < minPrediction {
// 				minPrediction = pred
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
