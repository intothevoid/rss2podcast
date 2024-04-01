package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Config struct {
	Subject        string `json:"subject"`
	Podcaster      string `json:"podcaster"`
	TtsUrl         string `json:"tts_url"`
	RssMaxArticles string `json:"rss_max_articles"`
	OllamaEndPoint string `json:"ollama_endpoint"`
	OllamaModel    string `json:"ollama_model"`
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// Handle post request with values podcaster, rssUrl, rssMaxArticles,
// ollama_endpoint, ollama_model, tts_url
// Use the values to update the config.yaml file and restart the podcaster
func ConfigureHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	// Check if the Content-Type is application/json
	// w.Header().Set("Content-Type", "application/json")
	// ct := r.Header.Get("Content-Type")
	// if ct != "" {
	// 	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	// 	if mediaType != "application/json" {
	// 		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
	// 		return
	// 	}
	// }

	var config Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := checkIfConfigExists(config.Subject, config.Podcaster,
		config.TtsUrl, config.RssMaxArticles, config.OllamaEndPoint,
		config.OllamaModel); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func checkIfConfigExists(subject string, podcaster string,
	ttsUrl string, rssMaxArticles string, ollamaEndPoint string,
	ollamaModel string) error {
	if podcaster == "" || subject == "" || ttsUrl == "" ||
		rssMaxArticles == "" || ollamaEndPoint == "" || ollamaModel == "" {
		return errors.New("invalid configuration")
	} else {
		return nil

	}
}
