package strategy

import (
	"sync"

	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
)

const (
	kRoundRobin         = "RoundRobin"
	kWeightedRoundRobin = "WeightedRoundRobin"
	kUnknown            = "Unknown"
)

type BalancingStrategy interface {
	Next([]*domain.Server) (*domain.Server, error)
}

var strategies map[string]func() BalancingStrategy

func init() {
	strategies[kRoundRobin] = func() BalancingStrategy {
		return &RoundRobin{
			mu:      sync.Mutex{},
			current: 0,
		}
	}

	strategies[kWeightedRoundRobin] = func() BalancingStrategy {
		return &WeightedRoundRobin{mu: sync.Mutex{}}
	}
}

// INFO: Based simple incrementable counter, just go to next server, current+1
type RoundRobin struct {
	mu      sync.Mutex
	current int
}

// Next implements BalancingStrategy.
func (r *RoundRobin) Next([]*domain.Server) (*domain.Server, error) {
	panic("unimplemented")
}

// INFO: RoundRobin with weight
type WeightedRoundRobin struct {
	count   []int
	current int
	mu      sync.Mutex
}

// Next implements BalancingStrategy.
func (w *WeightedRoundRobin) Next([]*domain.Server) (*domain.Server, error) {
	panic("unimplemented")
}
