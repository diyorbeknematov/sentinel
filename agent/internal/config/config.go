package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentName       string        `yaml:"agent_name"`
	AgentID         string        `yaml:"agent_id"`
	APIKey          string        `yaml:"api_key"`
	ServerURL       string        `yaml:"server_url"`
	AppLog          string        `yaml:"app_log"`
	NginxLog        string        `yaml:"nginx_log"`
	MetricsInterval time.Duration `yaml:"metrics_interval"`
	KafkaBrokers    []string      `yaml:"kafka_brokers"`
	KafkaTopic      string        `yaml:"kafka_topic"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.AgentName == "" {
		return fmt.Errorf("agent_name is required")
	}
	if c.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}
	if len(c.KafkaBrokers) == 0 {
		return fmt.Errorf("kafka_brokers is required")
	}
	if c.KafkaTopic == "" {
		return fmt.Errorf("kafka_topic is required")
	}
	if c.MetricsInterval <= 0 {
		c.MetricsInterval = 5 * time.Second
	}
	return nil
}

func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
