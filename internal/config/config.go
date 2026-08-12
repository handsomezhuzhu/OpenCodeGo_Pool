package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Scraper  ScraperConfig  `yaml:"scraper"`
	CPA      CPAConfig      `yaml:"cpa"`
}

type ServerConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type ScraperConfig struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

type CPAConfig struct {
	Endpoint     string `yaml:"endpoint"`
	BearerToken  string `yaml:"bearer_token"`
	ProviderName string `yaml:"provider_name"`
	BaseURL      string `yaml:"base_url"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Address:  ":8080",
			Password: "admin",
		},
		Database: DatabaseConfig{
			Path: "./data/pool.db",
		},
		Scraper: ScraperConfig{
			Interval: 5 * time.Minute,
			Timeout:  30 * time.Second,
		},
		CPA: CPAConfig{
			ProviderName: "OpenCode Go",
			BaseURL:      "https://opencode.ai/zen/go/v1",
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
