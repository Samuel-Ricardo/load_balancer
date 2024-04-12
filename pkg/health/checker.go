package health

import (
	"errors"
	"time"

	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	"github.com/pingcap/log"
)

type HealthChecker struct {
	servers []*domain.Server
	period  int
}

func NewChecker(_conf *domain.Config, servers []*domain.Server) (*HealthChecker, error) {
	if len(servers) == 0 {
		return nil, errors.New(`a server list expected, none provided or is empty list`)
	}

	return &HealthChecker{
		servers: servers,
	}, nil
}

func (hc *HealthChecker) start() {
	log.Info("Starting the health checker...")

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:
			for _, server := range hc.servers {
				go checkHealth(server)
			}
		}
	}
}
