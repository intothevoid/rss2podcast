package html

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Scrape retrieves the content of a web page specified by the given URL.
// It returns the extracted text from paragraph tags in the HTML document.
func Scrape(url string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // This allows redirects to be followed automatically
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Extract text from paragraph tags
	webContent := ""
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		webContent += s.Text()
	})

	// Optionally, extract text from other text-containing elements
	// doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
	//     fmt.Println(s.Text())
	// })

	return webContent
}
