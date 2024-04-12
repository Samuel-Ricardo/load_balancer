package health

import (
	"errors"

	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
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
