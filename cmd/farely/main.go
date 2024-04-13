package main

import (
	"flag"
	"net/http/httputil"
	"net/url"

	"github.com/Samuel-Ricardo/load_balancer/pkg/config"
	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	"github.com/Samuel-Ricardo/load_balancer/pkg/health"
	"github.com/Samuel-Ricardo/load_balancer/pkg/strategy"

	log "github.com/sirupsen/logrus"
)

var (
	port      = flag.Int("port", 8080, "Where to start farely")
	configure = flag.String("confi-path", "", "The config file to supply farely")
)

type Farely struct {
	Config     *config.Config
	ServerList map[string]*config.ServerList
}

func NewFarely(conf *config.Config) *Farely {
	serverMap := make(map[string]*config.ServerList, 0)

	for _, service := range conf.Services {

		servers := make([]*domain.Server, 0)

		for _, replica := range service.Replicas {

			url, err := url.Parse(replica.Url)
			if err != nil {
				log.Fatal("Could not parse URL")
			}

			proxy := httputil.NewSingleHostReverseProxy(url)
			servers = append(servers, &domain.Server{
				Url:      url,
				Proxy:    proxy,
				Metadata: replica.Metadata,
			})
		}

		checker, err := health.NewChecker(nil, servers)
		if err != nil {
			log.Fatal(err)
		}

		serverMap[service.Matcher] = &config.ServerList{
			Servers:  servers,
			Name:     service.Name,
			Strategy: strategy.LoadStrategy(service.Strategy),
			Hc:       checker,
		}
	}

	for _, sl := range serverMap {
		go sl.Hc.Start()
	}

	return &Farely{
		Config:     conf,
		ServerList: serverMap,
	}
}
