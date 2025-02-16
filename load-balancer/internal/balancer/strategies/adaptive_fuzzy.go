package strategies

// import (
// 	"math"
// 	"sync"

// 	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
// 	"github.com/tylerDurdenGolang/load-balancer/internal/errs"
// )

// type FuzzyBalancer struct {
// 	backends []*domain.Backend
// 	mutex    sync.RWMutex
// }

// func New(backends []*domain.Backend) *FuzzyBalancer {
// 	return &FuzzyBalancer{
// 		backends: backends,
// 	}
// }

// func (f *FuzzyBalancer) membership(val, low, medium, high float64) (float64, float64, float64) {
// 	lowM := math.Max(0, 1-(val-low)/(medium-low))
// 	mediumM := math.Max(0, 1-math.Abs(val-medium)/(high-medium))
// 	highM := math.Max(0, (val-medium)/(high-medium))
// 	return lowM, mediumM, highM
// }

// func (f *FuzzyBalancer) inferWeight(b *domain.Backend) float64 {
// 	cpu := b.Metrics.CPU
// 	lat := float64(b.Metrics.Latency.Milliseconds())

// 	// Фаззификация
// 	cpuLow, cpuMed, cpuHigh := f.membership(cpu, 0.2, 0.5, 0.8)
// 	latLow, latMed, latHigh := f.membership(lat, 50, 200, 500)

// 	// Правила вывода (Mamdani)
// 	rules := []float64{
// 		math.Min(cpuLow, latLow),   // Rule 1: Если CPU низкий И латентность низкая
// 		math.Min(cpuMed, latLow),   // Rule 2: Если CPU средний И латентность низкая
// 		math.Min(cpuHigh, latHigh), // Rule 3: Если CPU высокий И латентность высокая
// 	}

// 	// Дефаззификация (метод центра тяжести)
// 	numerator := rules[0]*0.2 + rules[1]*0.5 + rules[2]*0.8
// 	denominator := rules[0] + rules[1] + rules[2]

// 	if denominator == 0 {
// 		return 0.0
// 	}
// 	return numerator / denominator
// }

// func (f *FuzzyBalancer) GetBackend() (*domain.Backend, error) {
// 	f.mutex.RLock()
// 	defer f.mutex.RUnlock()

// 	maxWeight := -1.0
// 	var selected *domain.Backend

// 	for _, b := range f.backends {
// 		b.Lock()
// 		if b.Alive {
// 			weight := f.inferWeight(b)
// 			if weight > maxWeight {
// 				maxWeight = weight
// 				selected = b
// 			}
// 		}
// 		b.Unlock()
// 	}

// 	if selected == nil {
// 		return nil, errs.ErrNoHealthyBackends
// 	}
// 	return selected, nil
// }
