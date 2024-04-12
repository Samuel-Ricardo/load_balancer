package health

import (
	"errors"
	"net"
	"time"

	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	log "github.com/sirupsen/logrus"
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

func checkHealth(server *domain.Server) {
	_, err := net.DialTimeout("tcp", server.Url.Host, time.Second*5)
	if err != nil {

		log.Errorf("Could not connect to %s: %s, error: %s", server.Url.Host, server.Url.Path, err.Error())
		old := server.SetLiveness(false)

		if old {
			log.Warnf("Server %s is dead", server.Url.Host)
			log.Warnf("Transitioning server '%s' from Live to Unavailable state", server.Url.Host)
		}

		return
	}

	old := server.SetLiveness(true)

	if old {
		log.Warnf("Server %s is alive", server.Url.Host)
		log.Warnf("Transitioning server '%s' from Unavailable to Live state", server.Url.Host)
	}
}
