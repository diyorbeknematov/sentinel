package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AppLog   string `yaml:"app_log"`
	NginxLog string `yaml:"nginx_log"`
}

func Load() *Config {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("config.yaml not found:", err)
	}
	defer file.Close()

	cfg := &Config{}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal("failed to parse yaml:", err)
	}

	return cfg
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}

	return value
}
