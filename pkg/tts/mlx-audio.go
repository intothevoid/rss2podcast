package tts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type MLXAudioConverter struct {
	config *ConverterConfig
}

func NewMLXAudioConverter(config *ConverterConfig) *MLXAudioConverter {
	return &MLXAudioConverter{
		config: config,
	}
}

// ConvertToAudio sends a POST request to the MLX Audio TTS API and saves the response as an audio file.
func (c *MLXAudioConverter) ConvertToAudio(content string, fileName string) error {
	// Create form data
	formData := url.Values{}
	formData.Set("text", content)
	formData.Set("voice", c.config.Voice)
	formData.Set("speed", fmt.Sprintf("%.1f", c.config.Speed))

	// Create a new request
	req, err := http.NewRequest("POST", c.config.URL+"/tts", strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	// Download the audio file
	audioURL := c.config.URL + "/audio/" + response.Filename
	audioResp, err := http.Get(audioURL)
	if err != nil {
		return fmt.Errorf("error downloading audio file: %v", err)
	}
	defer audioResp.Body.Close()

	if audioResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(audioResp.Body)
		return fmt.Errorf("unexpected status code when downloading: %d, body: %s", audioResp.StatusCode, string(body))
	}

	// Create the output file
	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Copy the response body to the file
	_, err = io.Copy(outFile, audioResp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// PlayAudio sends a request to play the audio file directly from the server
func (c *MLXAudioConverter) PlayAudio(fileName string) error {
	formData := url.Values{}
	formData.Set("filename", filepath.Base(fileName))

	req, err := http.NewRequest("POST", c.config.URL+"/play", strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating play request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending play request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// StopAudio sends a request to stop any currently playing audio
func (c *MLXAudioConverter) StopAudio() error {
	req, err := http.NewRequest("POST", c.config.URL+"/stop", nil)
	if err != nil {
		return fmt.Errorf("error creating stop request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending stop request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// OpenOutputFolder sends a request to open the output folder in the system's file explorer
func (c *MLXAudioConverter) OpenOutputFolder() error {
	req, err := http.NewRequest("POST", c.config.URL+"/open_output_folder", nil)
	if err != nil {
		return fmt.Errorf("error creating open folder request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending open folder request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
