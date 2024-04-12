package strategy

import "github.com/Samuel-Ricardo/load_balancer/pkg/domain"

const (
	kRoundRobin         = "RoundRobin"
	kWeightedRoundRobin = "WeightedRoundRobin"
	kUnknown            = "Unknown"
)

type BalancingStrategy interface {
	Next([]*domain.Server) (*domain.Server, error)
}

var stratefies map[string]func() BalancingStrategy
