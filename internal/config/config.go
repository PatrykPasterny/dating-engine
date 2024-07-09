package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		URI        string `yaml:"uri"`
		Name       string `yaml:"name"`
		Collection string `yaml:"collection"`
	} `yaml:"database"`
	Redis struct {
		URI      string `yaml:"uri"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	}
	PageSize int64 `yaml:"pageSize"`
}

func GetConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed opening config file: %w", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("failed closing config file: %v", err)
		}
	}()

	var cfg Config

	decoder := yaml.NewDecoder(f)

	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed opening config file: %w", err)
	}

	return &cfg, err
}
