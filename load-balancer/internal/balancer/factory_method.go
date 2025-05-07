package balancer

import (
	"github.com/tylerDurdenGolang/load-balancer/internal/balancer/strategies"
	"github.com/tylerDurdenGolang/load-balancer/internal/domain"
)

func NewStrategy(algorithms string, backends []*domain.Backend) (IBalancer, error) {
	switch algorithms {
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
