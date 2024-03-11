package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ollama struct {
	endpoint string
	model    string
}

func NewOllama(endpoint string, model string) LLM {
	return &ollama{
		endpoint: endpoint,
		model:    model,
	}
}

// Define the request payload structure
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Format string `json:"format"`
	Stream bool   `json:"stream"`
}

// Update the GenerateResponse struct to include all fields from the JSON response
type GenerateResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// SendRequest sends a pre-crafted prompt to the Ollama API and returns the reply.
func (o *ollama) SendRequest(prompt string) (string, error) {
	// Create the request payload
	payload := GenerateRequest{
		Model:  o.model, // Using the model field from the ollama struct
		Prompt: prompt,
		Stream: false,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error marshaling payload: %v", err)
	}

	// Create a new HTTP request using the endpoint field
	req, err := http.NewRequest("POST", o.endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Create a JSON decoder for the response body
	decoder := json.NewDecoder(resp.Body)

	// Loop through the JSON stream
	var response GenerateResponse
	if err := decoder.Decode(&response); err != io.EOF {
		if err != nil {
			return "", fmt.Errorf("Error decoding JSON stream: %v", err)
		}
	}

	// Process each response object
	fmt.Printf("Model: %s, CreatedAt: %s, Response: %s, Done: %t\n", response.Model, response.CreatedAt, response.Response, response.Done)

	// Since we're processing a stream, we don't return a single response string.
	// You might want to aggregate responses or handle them differently based on your application's needs.
	return response.Response, nil
}
