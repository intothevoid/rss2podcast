package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/intothevoid/rss2podcast/internal/config"
	"github.com/intothevoid/rss2podcast/internal/io"
	"github.com/intothevoid/rss2podcast/internal/store"
	"github.com/intothevoid/rss2podcast/pkg/fileutil"
	"github.com/intothevoid/rss2podcast/pkg/llm"
	"github.com/intothevoid/rss2podcast/pkg/podcast"
	"github.com/intothevoid/rss2podcast/pkg/rss"
	"github.com/intothevoid/rss2podcast/pkg/tts"
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
	converter := tts.NewConverter(cfg.TTS.URL)
	writer := io.NewJsonWriter(store)

	// Check command line arguments
	noParse := false
	noConvert := false
	noConnectionCheck := false

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--no-parse":
			noParse = true
		case "--no-convert":
			noConvert = true
		case "--no-connection-check":
			noConnectionCheck = true
		}
	}

	// Check ollama connection
	if !noConnectionCheck {
		log.Println("Checking connection to Ollama...")
		err = ollama.CheckConnection()
		if err != nil {
			log.Fatal(err)
		}
	}

	// podcast filename
	// get timestamp as string in format yymmhh_hhmm
	ts := time.Now().Local().Format("2006_01_02_1504")
	podcast_fname_txt := fmt.Sprintf("%s_summary_%s.txt", cfg.Podcast.Subject, ts)
	podcast_fname_wav := fmt.Sprintf("%s_summary_%s.wav", cfg.Podcast.Subject, ts)

	tmpStore, err := writer.ReadStore()
	if err != nil {
		log.Println("No articles found, starting from scratch")
	} else {
		store = tmpStore
	}

	// Generate podcast introduction
	introduction, err := podcast.GenerateIntroduction(cfg.Podcast.Subject, cfg.Podcast.Podcaster)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating podcast introduction...")
	fileutil.AppendStringToFile(podcast_fname_txt, introduction)

	// Parse RSS feed
	if !noParse {
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
	}

	// Scrape all URLs and populate HTML content
	store.PopulateHtmlContent()

	// Write store to JSON
	writer.WriteStore(store)

	// Summarize articles
	for _, item := range store.GetData() {
		log.Printf("Summarizing article - %s", item.Title)
		summary, err := podcast.GenerateSummary(item.Title, item.Description, item.HtmlContent)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Done.")
		fileutil.AppendStringToFile(podcast_fname_txt, summary)
	}

	// Convert podcast text to audio
	if !noConvert {
		log.Println("Converting podcast text to audio...")
		fileContent, err := fileutil.ReadFileContent(podcast_fname_txt)
		if err != nil {
			log.Fatal(err)
		}
		converter.ConvertToAudio(fileContent, podcast_fname_wav)
	}
}
