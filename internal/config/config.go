package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RSS struct {
		URL         string   `yaml:"url"`
		MaxArticles int      `yaml:"max_articles"`
		Filters     []string `yaml:"filters"`
	} `yaml:"rss"`
	Ollama struct {
		EndPoint string `yaml:"end_point"`
		Model    string `yaml:"model"`
	} `yaml:"ollama"`
	Podcast struct {
		Subject   string `yaml:"subject"`
		Podcaster string `yaml:"podcaster"`
	} `yaml:"podcast"`
	TTS struct {
		Engine string `yaml:"engine"` // "coqui", "kokoro", or "mlx"
		Coqui  struct {
			URL string `yaml:"url"`
		} `yaml:"coqui"`
		Kokoro struct {
			URL    string  `yaml:"url"`
			Voice  string  `yaml:"voice"`
			Speed  float64 `yaml:"speed"`
			Format string  `yaml:"format"`
		} `yaml:"kokoro"`
		MLX struct {
			URL    string  `yaml:"url"`
			Voice  string  `yaml:"voice"`
			Speed  float64 `yaml:"speed"`
			Format string  `yaml:"format"`
		} `yaml:"mlx"`
	} `yaml:"tts"`
}

// LoadConfig loads the configuration from a YAML file.
// If the environment variable RSS2PODCAST_CONFIG is set, it will use the value as the path to the configuration file.
// Otherwise, it will default to "config.yaml" in the current directory.
// It returns a pointer to the loaded Config struct and any error encountered during the loading process.
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

// WriteConfig writes the provided configuration to a file.
// The configuration is encoded as YAML and saved to the file specified by the RSS2PODCAST_CONFIG environment variable.
// If the environment variable is not set, the configuration is saved to a file named "config.yaml" in the current directory.
// The function returns an error if there was a problem creating or writing to the file.
func WriteConfig(config *Config) error {
	configPath := os.Getenv("RSS2PODCAST_CONFIG")
	if configPath == "" {
		configPath = "config.yaml" // default to current directory
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}
