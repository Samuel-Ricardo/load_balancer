package domain

import (
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
