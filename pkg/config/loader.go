package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

func LoadConfig(reader io.Reader) (*Config, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	conf := Config{}
	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
