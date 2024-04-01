package handler

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	rss2podcast "github.com/intothevoid/rss2podcast/internal/app"
)

// GenerateHandler handles the /generate/{topic} endpoint.
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	// Get the topic from the request URL parameters
	vars := mux.Vars(r)
	topic := vars["topic"]
	app := rss2podcast.NewRSS2Podcast()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Check if the topic is empty
	if topic == "" {
		http.Error(w, "Invalid topic", http.StatusBadRequest)
		return
	}

	// Set the topic for the podcast
	app.SetTopic(topic)

	// Generate the .mp3 file using the topic
	generatedMp3, err := app.Run()
	if err != nil {
		http.Error(w, "Failed to generate .mp3 file", http.StatusInternalServerError)
		return
	}

	// Serve the generated .mp3 file
	http.ServeFile(w, r, generatedMp3)
}

// GenerateHandler handles the /generate/{topic} endpoint.
func GenerateHandlerTest(w http.ResponseWriter, r *http.Request) {
	// Get the topic from the request URL parameters
	vars := mux.Vars(r)
	topic := vars["topic"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Check if the topic is empty
	if topic == "" {
		http.Error(w, "Invalid topic", http.StatusBadRequest)
		return
	}

	testMp3 := "../../sample/News_summary_2024_03_17_1938.mp3"

	// check if the file exists
	if _, err := os.Stat(testMp3); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the generated .mp3 file
	http.ServeFile(w, r, testMp3)
}
