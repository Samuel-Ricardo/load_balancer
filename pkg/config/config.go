package config

import (
	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	"github.com/Samuel-Ricardo/load_balancer/pkg/health"
	"github.com/Samuel-Ricardo/load_balancer/pkg/strategy"
)

type Config struct {
	Services []domain.Service `yaml:"services"`
	Strategy string           `yaml:"strategy"`
}

type ServerList struct {
	Servers  []*domain.Server
	Name     string
	Strategy strategy.BalancingStrategy
	Hc       *health.HealthChecker
}
