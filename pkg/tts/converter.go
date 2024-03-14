package tts

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Converter struct {
	baseURL string
}

func NewConverter(baseURL string) *Converter {
	return &Converter{
		baseURL: baseURL,
	}
}

// ConvertToAudio sends a GET request with the specified content as a query parameter.
func (c *Converter) ConvertToAudio(content string, fileName string) error {
	// Create the request body
	params := url.Values{}
	params.Add("text", content)
	params.Add("speaker_id", "")
	params.Add("style_wav", "")
	params.Add("language_id", "")
	requestBody := strings.NewReader(params.Encode())

	// Create a new request using http
	req, err := http.NewRequest("POST", c.baseURL, requestBody)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return err
	}

	// Set the headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "http://localhost:5002/")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return err
	}

	// Save the response body to a file
	err = os.WriteFile(fileName, body, fs.FileMode(0644))
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
		return err
	}

	return nil
}
