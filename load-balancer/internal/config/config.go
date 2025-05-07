package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr          string        `yaml:"listen_addr"`
	Algorithm           string        `yaml:"algorithm"`
	WorkerPort          int           `yaml:"worker_port"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
