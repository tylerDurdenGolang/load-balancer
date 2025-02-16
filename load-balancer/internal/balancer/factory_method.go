package balancer

import (
	strategies "github.com/tylerDurdenGolang/load-balancer/internal/balancer/strategies"
	"github.com/tylerDurdenGolang/load-balancer/internal/config"
	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
)

func NewStrategy(config *config.Config) (IBalancer, error) {
	backends := make([]*domain.Backend, len(config.Backends))
	for i, backend := range config.Backends {
		backends[i] = domain.NewBackend(backend)
	}
	switch config.Algorithm {
	// case "fuzzy":
	// 	return adaptive_fuzzy.New(config.Backends), nil
	// case "queuing":
	// 	return queuing_theory.New(config.Backends, config.Lambda), nil
	// case "exponential":
	// 	return predictive.New(config.Backends, config.Alpha), nil
	// case "weighted":
	// 	return strategies.NewWeightedRandom(config.Backends, config.Alphas...), nil
	default:
		return strategies.NewRoundRobin(backends), nil
	}
}
