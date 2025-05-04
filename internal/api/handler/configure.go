package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/intothevoid/rss2podcast/internal/config"
)

type ConfigWebSvc struct {
	Subject        string `json:"subject"`
	Podcaster      string `json:"podcaster"`
	RssMaxArticles string `json:"rss_max_articles"`
	OllamaEndPoint string `json:"ollama_endpoint"`
	OllamaModel    string `json:"ollama_model"`
	TtsEngine      string `json:"tts_engine"`
	CoquiUrl       string `json:"coqui_url"`
	KokoroUrl      string `json:"kokoro_url"`
	KokoroVoice    string `json:"kokoro_voice"`
	KokoroSpeed    string `json:"kokoro_speed"`
	KokoroFormat   string `json:"kokoro_format"`
	MLXUrl         string `json:"mlx_url"`
	MLXVoice       string `json:"mlx_voice"`
	MLXSpeed       string `json:"mlx_speed"`
	MLXFormat      string `json:"mlx_format"`
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func checkIfConfigExists(subject string, podcaster string,
	rssMaxArticles string, ollamaEndPoint string,
	ollamaModel string) error {
	// Check required fields
	if subject == "" {
		return errors.New("subject is required")
	}
	if podcaster == "" {
		return errors.New("podcaster is required")
	}
	if rssMaxArticles == "" {
		return errors.New("rss_max_articles is required")
	}
	if ollamaEndPoint == "" {
		return errors.New("ollama_endpoint is required")
	}
	if ollamaModel == "" {
		return errors.New("ollama_model is required")
	}
	return nil
}

// Handle post request with values podcaster, rssUrl, rssMaxArticles,
// ollama_endpoint, ollama_model, tts_url
// Use the values to update the config.yaml file and restart the podcaster
func ConfigureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received configuration request")
	log.Printf("Request Method: %s", r.Method)
	log.Printf("Request Headers: %v", r.Header)

	// Enable CORS
	enableCORS(&w)

	// Handle preflight request
	if r.Method == "OPTIONS" {
		log.Println("Handling OPTIONS request")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and log the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received request body: %s", string(body))

	// Create a new reader with the read body
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var confIncoming ConfigWebSvc
	err = json.NewDecoder(r.Body).Decode(&confIncoming)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received configuration: %+v", confIncoming)

	if err := checkIfConfigExists(confIncoming.Subject, confIncoming.Podcaster,
		confIncoming.RssMaxArticles, confIncoming.OllamaEndPoint,
		confIncoming.OllamaModel); err != nil {
		log.Printf("Invalid configuration: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conf, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		http.Error(w, "Failed to load config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the config
	conf.Podcast.Subject = confIncoming.Subject
	conf.Podcast.Podcaster = confIncoming.Podcaster
	conf.TTS.Engine = confIncoming.TtsEngine
	conf.TTS.Coqui.URL = confIncoming.CoquiUrl
	conf.TTS.Kokoro.URL = confIncoming.KokoroUrl
	conf.TTS.Kokoro.Voice = confIncoming.KokoroVoice

	// Convert speed string to float64 for Kokoro
	if speed, err := strconv.ParseFloat(confIncoming.KokoroSpeed, 64); err == nil {
		conf.TTS.Kokoro.Speed = speed
	} else {
		log.Printf("Warning: Invalid Kokoro speed value '%s', using default", confIncoming.KokoroSpeed)
	}

	conf.TTS.Kokoro.Format = confIncoming.KokoroFormat

	// Update MLX configuration
	conf.TTS.MLX.URL = confIncoming.MLXUrl
	conf.TTS.MLX.Voice = confIncoming.MLXVoice

	// Convert speed string to float64 for MLX
	if speed, err := strconv.ParseFloat(confIncoming.MLXSpeed, 64); err == nil {
		conf.TTS.MLX.Speed = speed
	} else {
		log.Printf("Warning: Invalid MLX speed value '%s', using default", confIncoming.MLXSpeed)
	}

	conf.TTS.MLX.Format = confIncoming.MLXFormat

	maxArticles, err := strconv.Atoi(confIncoming.RssMaxArticles)
	if err != nil {
		log.Printf("Invalid max articles value: %v", err)
		http.Error(w, "Invalid max articles value: "+err.Error(), http.StatusBadRequest)
		return
	}

	conf.RSS.MaxArticles = maxArticles
	conf.Ollama.EndPoint = confIncoming.OllamaEndPoint
	conf.Ollama.Model = confIncoming.OllamaModel

	log.Printf("Saving configuration: %+v", conf)
	err = config.WriteConfig(conf)
	if err != nil {
		log.Printf("Failed to write config: %v", err)
		http.Error(w, "Failed to write config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Configuration saved successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Configuration saved successfully"})
}
