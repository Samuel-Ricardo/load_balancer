package domain

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Replica struct {
	Metadata map[string]string `yaml:"metadata"`
	Url      string            `yaml:"url"`
}

// TODO: Improve Matcher with something like Regex or Subsdomain
type Service struct {
	Name     string    `yaml:"name"`
	Matcher  string    `yaml:"matcher"`
	Strategy string    `yaml:"strategy"`
	Replicas []Replica `yaml:"replicas"`
}

type Config struct {
	Strategy string    `yaml:"strategy"`
	Services []Service `yaml:"services"`
}

type Server struct {
	Url      *url.URL
	Proxy    *httputil.ReverseProxy
	Metadata map[string]string
	mu       sync.RWMutex
	alive    bool
}

func (s *Server) Forward(res http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(res, req)
}

func (s *Server) GetMetaOrDefault(key, def string) string {
	v, ok := s.Metadata[key]

	if !ok {
		return def
	}

	return v
}

func (s *Server) SetLiveness(value bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	old := s.alive
	s.alive = value

	return old
}
