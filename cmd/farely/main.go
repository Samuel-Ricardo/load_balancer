package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Samuel-Ricardo/load_balancer/pkg/config"
	"github.com/Samuel-Ricardo/load_balancer/pkg/domain"
	"github.com/Samuel-Ricardo/load_balancer/pkg/health"
	"github.com/Samuel-Ricardo/load_balancer/pkg/strategy"

	log "github.com/sirupsen/logrus"
)

var (
	port       = flag.Int("port", 8080, "Where to start farely")
	configFile = flag.String("config-path", "", "The config file to supply farely")
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

func (f *Farely) findServiceList(reqPath string) (*config.ServerList, error) {
	log.Infof("Trying to find matcher for request: %s", reqPath)

	for matcher, s := range f.ServerList {
		if strings.HasPrefix(reqPath, matcher) {
			log.Infof("Found service '%s' matching the request", s.Name)
			return s, nil
		}
	}

	return nil, fmt.Errorf("could not find a matcher for url: '%s'", reqPath)
}

func (f *Farely) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Infof("Received new request: url='%s'", req.Host)

	sl, err := f.findServiceList(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	next, err := sl.Strategy.Next(sl.Servers)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("Fowarding to the server='%s'", next.Url.Host)
	next.Forward(res, req)
}

func main() {
	flag.Parse()

	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conf, err := config.LoadConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	farely := NewFarely(conf)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: farely,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
