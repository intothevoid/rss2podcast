package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Ollama struct {
	endpoint string
}

func NewOllama(endpoint string) *Ollama {
	return &Ollama{
		endpoint: endpoint,
	}
}

func (o *Ollama) Summarize(text string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(o.endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["summary"], nil
}
