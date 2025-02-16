package strategies

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/tylerDurdenGolang/load-balancer/internal/errs"
)

type Metrics struct {
	CPU       float64       // от 0 до 1 (или 0..100 если хотите)
	MemUsage  float64       // аналогично
	Latency   time.Duration // среднее время ответа
	ErrorRate float64       // от 0 до 1
}

// Backend содержит информацию о сервере и его метриках
type Backend struct {
	Host  string
	Alive bool
	mux   sync.RWMutex

	// Модель
	Metrics Metrics

	// Расчётная оценка (score) и вес (weight)
	Score  float64
	Weight float64
}

// WeightedRandom — структура с набором бэкендов и параметрами для вычисления score
type WeightedRandom struct {
	backends []*Backend

	// Коэффициенты для расчёта score
	alphaCPU       float64
	alphaMem       float64
	alphaLatency   float64
	alphaErrorRate float64

	randMutex sync.Mutex // для защиты rand.Seed
}

// NewWeightedRandom инициализирует балансировщик
func NewWeightedRandom(backends []*Backend,
	alphaCPU, alphaMem, alphaLatency, alphaError float64) *WeightedRandom {

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	wr := &WeightedRandom{
		backends:       backends,
		alphaCPU:       alphaCPU,
		alphaMem:       alphaMem,
		alphaLatency:   alphaLatency,
		alphaErrorRate: alphaError,
	}

	// Можно один раз пересчитать, или пересчитывать периодически
	wr.RecalcScoresAndWeights()
	return wr
}

// RecalcScoresAndWeights пересчитывает score и weight для каждого бэкенда.
// Вызывается периодически или после обновления метрик.
func (wr *WeightedRandom) RecalcScoresAndWeights() {
	for _, b := range wr.backends {
		b.mux.Lock()

		if !b.Alive {
			// Если бэкенд не живой, можно выставлять weight=0, score=∞
			b.Score = math.Inf(1)
			b.Weight = 0
			b.mux.Unlock()
			continue
		}

		// Вычислим score как линейную комбинацию
		cpu := b.Metrics.CPU
		mem := b.Metrics.MemUsage
		lat := float64(b.Metrics.Latency.Milliseconds()) // переведём в ms
		errRate := b.Metrics.ErrorRate

		// Пример: score = α1*CPU + α2*Mem + α3*Latency + α4*Err
		s := wr.alphaCPU*cpu + wr.alphaMem*mem + wr.alphaLatency*lat + wr.alphaErrorRate*errRate
		if s <= 0 {
			// минимальная защита от деления на ноль
			s = 0.0001
		}

		b.Score = s
		b.Weight = 1.0 / s // обратная пропорция

		b.mux.Unlock()
	}
}

// GetBackend выбирает бэкенд случайным образом, пропорционально весам
func (wr *WeightedRandom) GetBackend() (string, error) {
	// Суммируем веса живых бэкендов
	var totalWeight float64
	var aliveBackends []*Backend

	for _, b := range wr.backends {
		b.mux.RLock()
		w := b.Weight
		alive := b.Alive && w > 0
		b.mux.RUnlock()

		if alive {
			aliveBackends = append(aliveBackends, b)
			totalWeight += w
		}
	}

	if len(aliveBackends) == 0 {
		return "", errs.ErrNoHealthyBackends
	}

	// Генерируем случайное число [0..totalWeight]
	wr.randMutex.Lock() // На случай, если много горутин
	r := rand.Float64() * totalWeight
	wr.randMutex.Unlock()

	// Ищем, какой бэкенд выпал
	for _, b := range aliveBackends {
		b.mux.RLock()
		w := b.Weight
		b.mux.RUnlock()

		if r < w {
			return b.Host, nil
		}
		r -= w
	}

	// Теоретически не должно сюда доходить, но на всякий случай:
	last := aliveBackends[len(aliveBackends)-1]
	return last.Host, nil
}
