package handler

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

// GenerateHandler handles the /generate/{topic} endpoint.
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	// Get the topic from the request URL parameters
	vars := mux.Vars(r)
	topic := vars["topic"]

	// Check if the topic is empty
	if topic == "" {
		http.Error(w, "Invalid topic", http.StatusBadRequest)
		return
	}

	// Generate the .mp3 file using the topic
	err := generateMP3(topic)
	if err != nil {
		http.Error(w, "Failed to generate .mp3 file", http.StatusInternalServerError)
		return
	}

	// Serve the generated .mp3 file
	http.ServeFile(w, r, "generated.mp3")
}

// generateMP3 generates a .mp3 file using the given topic.
func generateMP3(topic string) error {
	// TODO: Implement the logic to generate the .mp3 file using the topic.
	// You can use any external libraries or APIs to generate the .mp3 file.
	// Here's an example using the `say` command on macOS:
	cmd := exec.Command("say", "-o", "generated.mp3", topic)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate .mp3 file: %w", err)
	}

	return nil
}
