package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary config file
	testData := `
rss:
  url: http://example.com/feed
`
	_, err = tempFile.WriteString(testData)
	if err != nil {
		t.Fatalf("Failed to write test data to config file: %v", err)
	}
	tempFile.Close()

	// Set the RSS2PODCAST_CONFIG environment variable to the temporary config file
	os.Setenv("RSS2PODCAST_CONFIG", tempFile.Name())

	// Call the LoadConfig function
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the loaded config values
	expectedURL := "http://example.com/feed"
	if cfg.RSS.URL != expectedURL {
		t.Errorf("Expected RSS URL to be %q, got %q", expectedURL, cfg.RSS.URL)
	}
}
