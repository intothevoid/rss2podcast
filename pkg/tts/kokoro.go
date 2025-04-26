package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type KokoroConverter struct {
	config *ConverterConfig
}

func NewKokoroConverter(config *ConverterConfig) *KokoroConverter {
	return &KokoroConverter{
		config: config,
	}
}

type KokoroRequest struct {
	Model              string  `json:"model"`
	Input              string  `json:"input"`
	Voice              string  `json:"voice"`
	ResponseFormat     string  `json:"response_format"`
	Speed              float64 `json:"speed"`
	Stream             bool    `json:"stream"`
	ReturnDownloadLink bool    `json:"return_download_link"`
}

// ConvertToAudio sends a POST request to the Kokoro TTS API and saves the response as an audio file.
func (c *KokoroConverter) ConvertToAudio(content string, fileName string) error {
	// Create the request payload
	payload := KokoroRequest{
		Model:              "kokoro",
		Input:              content,
		Voice:              c.config.Voice,
		ResponseFormat:     c.config.Format,
		Speed:              c.config.Speed,
		Stream:             false,
		ReturnDownloadLink: true,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", c.config.URL+"/v1/audio/speech", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "audio/mpeg")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Get the download link from headers if available
	downloadPath := resp.Header.Get("X-Download-Path")
	if downloadPath != "" {
		// Download the file from the provided path
		downloadURL := c.config.URL + downloadPath
		downloadResp, err := http.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("error downloading audio file: %v", err)
		}
		defer downloadResp.Body.Close()

		// Create the output file
		outFile, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}
		defer outFile.Close()

		// Copy the response body to the file
		_, err = io.Copy(outFile, downloadResp.Body)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	} else {
		// If no download link, save the response body directly
		outFile, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
