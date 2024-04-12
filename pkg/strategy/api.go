package strategy

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	"github.com/pingcap/log"
	"golang.org/x/crypto/openpgp/errors"
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

// NOTE: implements Round Robin algorithm for BalancingStrategy interface
func (r *RoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	seen := 0
	var picked *domain.Server

	for seen < len(servers) {
		picked = servers[r.current]
		r.current = (r.current + 1) % len(servers)

		if picked.IsAlive() {
			break
		}

		seen++
	}

	if picked == nil || seen == len(servers) {
		log.Error("All servers are dead")
		return nil, errors.New(fmt.Sprintf("Checked all the '%d' servers, none of them are alive", seen))
	}

	log.Infof("Strategy picked server: %s", picked.Url.Host)
	return picked, nil
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
