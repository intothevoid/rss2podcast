package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RSS struct {
		URL string `yaml:"url"`
	} `yaml:"rss"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
