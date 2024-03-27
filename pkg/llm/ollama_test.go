package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendRequest(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method and content type
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read the request body
		var payload GenerateRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		// Verify the payload values
		expectedModel := "test-model"
		expectedPrompt := "test-prompt"
		if payload.Model != expectedModel {
			t.Errorf("Expected model to be %s, got %s", expectedModel, payload.Model)
		}
		if payload.Prompt != expectedPrompt {
			t.Errorf("Expected prompt to be %s, got %s", expectedPrompt, payload.Prompt)
		}

		// Create a mock response
		response := GenerateResponse{
			Model:     "test-model",
			CreatedAt: time.Now(),
			Response:  "test-response",
			Done:      true,
		}

		// Encode the response as JSON and write it to the response writer
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Errorf("Error encoding response: %v", err)
		}
	}))
	defer server.Close()

	// Create an instance of the ollama struct
	o := &ollama{
		model:    "test-model",
		endpoint: server.URL,
	}

	// Call the SendRequest method
	prompt := "test-prompt"
	response, err := o.SendRequest(prompt)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}

	// Verify the response value
	expectedResponse := "test-response"
	if response != expectedResponse {
		t.Errorf("Expected response to be %s, got %s", expectedResponse, response)
	}
}
