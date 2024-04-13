package strategy

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	log "github.com/sirupsen/logrus"
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
	strategies = make(map[string]func() BalancingStrategy, 0)
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
		return nil, fmt.Errorf("checked all the '%d' servers, none of them are alive", seen)
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

// NOTE: implements Weighted Round Robin algorithm for BalancingStrategy interface
func (w *WeightedRoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.count == nil {
		w.count = make([]int, len(servers))
		w.current = 0
	}

	seen := 0
	var picked *domain.Server

	for seen < len(servers) {

		picked = servers[w.current]
		capacity := picked.GetMetaOrDefaultInt("weight", 1)

		if !picked.IsAlive() {
			seen++
			w.count[w.current] = 0
			w.current = (w.current + 1) % len(servers)
			continue
		}

		if w.count[w.current] <= capacity {

			w.count[w.current]++
			log.Infof("Strategy picked server: %s", picked.Url.Host)

			return picked, nil
		}

		w.count[w.current] = 0
		w.current = (w.current + 1) % len(servers)
	}

	if picked == nil || seen == len(servers) {
		log.Error("All servers are dead")
		return nil, errors.New(fmt.Sprintf("Checked all the '%d' servers, none of them is available", seen))
	}

	return picked, nil
}

func LoadStrategy(name string) BalancingStrategy {
	st, ok := strategies[name]
	if !ok {
		log.Warnf("Strategy with name '%s' not found, falling back to Round Robin Strategy", name)
		return strategies[kRoundRobin]()
	}
	log.Infof("Picked Strategy '%s'", name)
	return st()
}
