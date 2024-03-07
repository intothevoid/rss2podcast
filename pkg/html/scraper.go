package html

import (
	"io"
	"log"
	"net/http"
)

// Function which accepts a URL and returns the HTML content of the page
func Scrape(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.Reader(resp.Body))
	return string(body)
}
