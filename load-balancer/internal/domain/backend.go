package domain

import (
	"sync"
	"time"
)

// Backend представляет сервер, на который балансировщик направляет запросы.
type Backend struct {
	host    string
	alive   bool
	mux     sync.RWMutex
	metrics *Metrics
}

// Metrics содержит метрики производительности бэкенда.
type Metrics struct {
	cpuUsage        float64       // Загрузка CPU (0.0 - 1.0)
	memoryUsage     float64       // Использование памяти (0.0 - 1.0)
	activeRequests  int           // Текущее количество активных запросов
	avgResponseTime time.Duration // Среднее время ответа (скользящее среднее)
	errorRate       float64       // Доля ошибок (0.0 - 1.0)
	lastUpdated     time.Time     // Время последнего обновления метрик
}

// NewBackend создает новый экземпляр бэкенда.
func NewBackend(host string) *Backend {
	return &Backend{
		host:    host,
		alive:   true, // По умолчанию считается живым
		metrics: &Metrics{},
	}
}

// Host возвращает адрес бэкенда.
func (b *Backend) Host() string {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.host
}

// SetAlive устанавливает статус доступности бэкенда.
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.alive = alive
}

// IsAlive возвращает текущий статус доступности бэкенда.
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}

// UpdateMetrics обновляет метрики бэкенда.
func (b *Backend) UpdateMetrics(cpu, memory, errorRate float64, responseTime time.Duration) {
	b.mux.Lock()
	defer b.mux.Unlock()

	// Экспоненциальное сглаживание для среднего времени ответа
	alpha := 0.2
	if b.metrics.avgResponseTime == 0 {
		b.metrics.avgResponseTime = responseTime
	} else {
		b.metrics.avgResponseTime = time.Duration(
			(1-alpha)*float64(b.metrics.avgResponseTime) + alpha*float64(responseTime),
		)
	}

	b.metrics.cpuUsage = cpu
	b.metrics.memoryUsage = memory
	b.metrics.errorRate = errorRate
	b.metrics.lastUpdated = time.Now()
}

// Metrics возвращает текущие метрики бэкенда.
func (b *Backend) Metrics() Metrics {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return *b.metrics
}

// IncrementRequests увеличивает счетчик активных запросов.
func (b *Backend) IncrementRequests() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.metrics.activeRequests++
}

// DecrementRequests уменьшает счетчик активных запросов.
func (b *Backend) DecrementRequests() {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.metrics.activeRequests > 0 {
		b.metrics.activeRequests--
	}
}
