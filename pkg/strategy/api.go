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

// INFO: Based simple incrementable counter, just go to next server, current+1
type RoundRobin struct {
	mu      sync.Mutex
	current int
}

// INFO: RoundRobin with weight
type WeightedRoundRobin struct {
	count   []int
	current int
	mu      sync.Mutex
}
