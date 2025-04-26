package tts

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConverter_ConvertToAudio(t *testing.T) {
	// Create a test server to mock the HTTP request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request headers
		if r.Header.Get("User-Agent") != "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0" {
			t.Errorf("Unexpected User-Agent header: got %s, want Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("Accept") != "*/*" {
			t.Errorf("Unexpected Accept header: got %s, want */*", r.Header.Get("Accept"))
		}
		// Add more header verifications here...

		// Verify the request body
		body, err := io.ReadAll(io.Reader(r.Body))
		if err != nil {
			t.Fatalf("Error reading request body: %s", err)
		}
		expectedBody := "language_id=&speaker_id=&style_wav=&text=sample+content"
		if string(body) != expectedBody {
			t.Errorf("Unexpected request body: got %s, want %s", string(body), expectedBody)
		}

		// Send a mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock audio data"))
	}))
	defer server.Close()

	// Create config structs
	coquiConfig := &ConverterConfig{
		URL: server.URL,
	}

	kokoroConfig := &ConverterConfig{
		URL:    "",
		Voice:  "af_heart",
		Speed:  1.0,
		Format: "mp3",
	}

	// Create a Converter instance with the test server URL
	converter := &Converter{
		engine: "coqui",
		coqui:  NewCoquiConverter(coquiConfig),
		kokoro: NewKokoroConverter(kokoroConfig),
	}

	// Call the method being tested
	err := converter.ConvertToAudio("sample content", "test.wav")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Verify the file was written correctly
	fileData, err := os.ReadFile("test.wav")
	defer os.Remove("test.wav")
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}
	expectedFileData := []byte("mock audio data")
	if !bytes.Equal(fileData, expectedFileData) {
		t.Errorf("Unexpected file data: got %s, want %s", string(fileData), string(expectedFileData))
	}
}
