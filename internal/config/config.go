package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RSS struct {
		URL         string `yaml:"url"`
		MaxArticles int    `yaml:"max_articles"`
	} `yaml:"rss"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("RSS2PODCAST_CONFIG")
	if configPath == "" {
		configPath = "config.yaml" // default to current directory
	}

	f, err := os.Open(configPath)
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
