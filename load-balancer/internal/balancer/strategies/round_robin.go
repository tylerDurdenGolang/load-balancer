package strategies

import (
	"sync"

	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
	"github.com/tylerDurdenGolang/load-balancer/internal/errs"
)

type RoundRobin struct {
	backends []*domain.Backend
	current  int
	mutex    sync.Mutex
}

func NewRoundRobin(b []*domain.Backend) *RoundRobin {
	return &RoundRobin{backends: b}
}

func (r *RoundRobin) GetBackend() (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	aliveBackends := make([]*domain.Backend, 0, len(r.backends))
	for _, b := range r.backends {
		isAlive := b.IsAlive()
		if isAlive {
			aliveBackends = append(aliveBackends, b)
		}
	}
	if len(aliveBackends) == 0 {
		return "", errs.ErrNoHealthyBackends
	}

	backend := aliveBackends[r.current%len(aliveBackends)]
	r.current = (r.current + 1) % len(aliveBackends)

	return backend.Host(), nil
}

func (r *RoundRobin) MarkBackendDown(backend string) {
	// Удаляем (или помечаем) backend как «упавший»
	// Для простоты — уберём из списка
	r.mutex.Lock()
	defer r.mutex.Unlock()

	newList := make([]*domain.Backend, 0, len(r.backends))
	for _, b := range r.backends {
		if b.Host() != backend {
			newList = append(newList, b)
		}
	}
	r.backends = newList
}

func (r *RoundRobin) MarkBackendUp(backend string) {
	// Если бэкенд отсутствует в списке, добавим его
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, b := range r.backends {
		if b.Host() == backend {
			return // уже есть
		}
	}
	r.backends = append(r.backends, domain.NewBackend(backend))
}

func (r *RoundRobin) GetAllBackends() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	result := make([]string, len(r.backends))
	for i, b := range r.backends {
		result[i] = b.Host()
	}
	return result
}

func (r *RoundRobin) ReplaceBackends(newList []*domain.Backend) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.backends = newList
	r.current = 0
}
