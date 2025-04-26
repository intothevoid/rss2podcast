package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/intothevoid/rss2podcast/internal/config"
)

type ConfigWebSvc struct {
	Subject        string `json:"subject"`
	Podcaster      string `json:"podcaster"`
	TtsUrl         string `json:"tts_url"`
	RssMaxArticles string `json:"rss_max_articles"`
	OllamaEndPoint string `json:"ollama_endpoint"`
	OllamaModel    string `json:"ollama_model"`
	TtsEngine      string `json:"tts_engine"`
	CoquiUrl       string `json:"coqui_url"`
	KokoroUrl      string `json:"kokoro_url"`
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

	var confIncoming ConfigWebSvc
	err := json.NewDecoder(r.Body).Decode(&confIncoming)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := checkIfConfigExists(confIncoming.Subject, confIncoming.Podcaster,
		confIncoming.TtsUrl, confIncoming.RssMaxArticles, confIncoming.OllamaEndPoint,
		confIncoming.OllamaModel); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		conf, err := config.LoadConfig()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the config
		conf.Podcast.Subject = confIncoming.Subject
		conf.Podcast.Podcaster = confIncoming.Podcaster
		conf.TTS.Engine = confIncoming.TtsEngine
		conf.TTS.Coqui.URL = confIncoming.CoquiUrl
		conf.TTS.Kokoro.URL = confIncoming.KokoroUrl

		maxArticles, err := strconv.Atoi(confIncoming.RssMaxArticles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		conf.RSS.MaxArticles = maxArticles
		conf.Ollama.EndPoint = confIncoming.OllamaEndPoint
		conf.Ollama.Model = confIncoming.OllamaModel
		config.WriteConfig(conf)
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
