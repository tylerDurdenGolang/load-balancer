package balancer

import "github.com/tylerDurdenGolang/load-balancer/internal/domain"

// IBalancer — интерфейс для выбора бэкенда.
// Можно дополнить методами для управления состоянием бэкендов.
type IBalancer interface {
	GetBackend() (string, error)
	MarkBackendDown(backend string)
	MarkBackendUp(backend string)
	GetAllBackends() []string
	ReplaceBackends(newList []*domain.Backend)
}
