package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig(strings.NewReader(`
services:
- name: "test service" 
  strategy: RoundRobin 
  matcher: "/api/v1"
  replicas:
  - url: localhost:8081
  - url: localhost:8082
strategy: RoundRobin
    `))
	if err != nil {
		t.Errorf("Failed to load config: %s", err)
	}

	if conf.Strategy != "RoundRobin" {
		t.Errorf("Failed to load config, strategy expexted 'RoundRobin': got '%s'", conf.Strategy)
	}

	if len(conf.Services) != 1 {
		t.Errorf("Failed to load config, services expexted 1: got '%d'", len(conf.Services))
	}

	if conf.Services[0].Matcher != "/api/v1" {
		t.Errorf("Failed to load config, service name expexted '/api/v1': got '%s'", conf.Services[0].Matcher)
	}

	if conf.Services[0].Name != "test service" {
		t.Errorf("Failed to load config, service name expexted 'test service': got '%s'", conf.Services[0].Name)
	}

	if conf.Services[0].Strategy != "RoundRobin" {
		t.Errorf("Failed to load config, service name expexted 'RoundRobin': got '%s'", conf.Services[0].Strategy)
	}

	if len(conf.Services[0].Replicas) != 2 {
		t.Errorf("Failed to load config, replicas expexted 2: got '%d'", len(conf.Services[0].Replicas))
	}

	if conf.Services[0].Replicas[0].Url != "localhost:8081" {
		t.Errorf("Failed to load config, service name expexted 'localhost:8081': got '%s'", conf.Services[0].Replicas[0])
	}

	if conf.Services[0].Replicas[1].Url != "localhost:8082" {
		t.Errorf("Failed to load config, service name expexted 'localhost:8082': got '%s'", conf.Services[0].Replicas[1])
	}
}
