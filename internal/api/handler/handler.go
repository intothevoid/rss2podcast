package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	rss2podcast "github.com/intothevoid/rss2podcast/internal/app"
)

// GenerateHandler handles the /generate/{topic} endpoint.
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	// Get the topic from the request URL parameters
	vars := mux.Vars(r)
	topic := vars["topic"]
	app := rss2podcast.NewRSS2Podcast()

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
