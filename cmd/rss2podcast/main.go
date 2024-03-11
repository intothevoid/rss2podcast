package main

import (
	"fmt"
	"log"
	"time"

	"github.com/intothevoid/rss2podcast/internal/config"
	"github.com/intothevoid/rss2podcast/internal/io"
	"github.com/intothevoid/rss2podcast/internal/store"
	"github.com/intothevoid/rss2podcast/pkg/fileutil"
	"github.com/intothevoid/rss2podcast/pkg/llm"
	"github.com/intothevoid/rss2podcast/pkg/podcast"
	"github.com/intothevoid/rss2podcast/pkg/rss"
)

func main() {
	// Get RSS feed URL from config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Set up dependencies
	rssParser := rss.NewParser()
	store := store.NewStore()
	ollama := llm.NewOllama(cfg.Ollama.EndPoint, cfg.Ollama.Model)
	podcast := podcast.NewPodcast(ollama)
	// converter := tts.NewConverter()
	writer := io.NewJsonWriter(store)

	// podcast filename
	// get timestamp as string in format yymmhh_hhmm
	ts := time.Now().Local().Format("2006_01_02_1504")
	podcast_fname := fmt.Sprintf("%s_summary_%s.txt", cfg.Podcast.Subject, ts)

	store, err = writer.ReadStore()
	if err != nil {
		log.Println("No articles found, starting from scratch")
	}

	// Generate podcast introduction
	introduction, err := podcast.GenerateIntroduction(cfg.Podcast.Subject, cfg.Podcast.Podcaster)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Generating podcast introduction...")
	fileutil.AppendStringToFile(podcast_fname, introduction)

	// Summarize articles
	for _, item := range store.GetData() {
		log.Printf("Summarizing article - %s", item.Title)
		summary, err := podcast.GenerateSummary(item.Title, item.Description, item.HtmlContent)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Done.")
		fileutil.AppendStringToFile(podcast_fname, summary)
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
}
