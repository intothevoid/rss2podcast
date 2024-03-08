package main

import (
	"log"

	"github.com/intothevoid/rss2podcast/internal/config"
	"github.com/intothevoid/rss2podcast/internal/io"
	"github.com/intothevoid/rss2podcast/internal/store"
	"github.com/intothevoid/rss2podcast/pkg/llm"
	"github.com/intothevoid/rss2podcast/pkg/rss"
	"github.com/intothevoid/rss2podcast/pkg/tts"
)

func main() {
	rssParser := rss.NewParser()
	store := store.NewStore()
	ollama := llm.NewOllama("http://localhost:8000")
	converter := tts.NewConverter()
	writer := io.NewJsonWriter(store)

	// Get RSS feed URL from config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Parse RSS feed
	items, _ := rssParser.ParseURL(cfg.RSS.URL)

	// Store top 'maxArticles' articles, default to 10
	for i, item := range items {
		if i >= cfg.RSS.MaxArticles {
			break
		}

		rssItem := rss.RSSItem{
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			HtmlContent: "",
		}
		store.Save(item.GUID, rssItem)
	}

	// Scrape all URLs and populate HTML content
	store.PopulateHtmlContent()

	// Write store to JSON
	writer.WriteStore(store)

	// Summarize and convert to audio
	for _, rssItem := range store.GetData() {
		summary, _ := ollama.Summarize(rssItem.HtmlContent)
		converter.ConvertToAudio(summary, rssItem.Title+".aiff")
	}
}
