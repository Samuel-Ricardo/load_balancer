package main

import (
	"flag"

	"github.com/Samuel-Ricardo/load_balancer/pkg/config"
)

var (
	port      = flag.Int("port", 8080, "Where to start farely")
	configure = flag.String("confi-path", "", "The config file to supply farely")
)

type Farely struct {
	Config     *config.Config
	ServerList map[string]*config.ServerList
}
