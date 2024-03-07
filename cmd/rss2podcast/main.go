package main

import (
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

	// Parse RSS feed
	items, _ := rssParser.ParseURL("http://example.com/rss")

	// Store top 10 articles
	for i, item := range items {
		if i >= 10 {
			break
		}
		store.Save(item.Title, item.Description)
	}

	// Summarize and convert to audio
	for title, description := range store.GetData() {
		summary, _ := ollama.Summarize(description)
		converter.ConvertToAudio(summary, title+".aiff")
	}
}
