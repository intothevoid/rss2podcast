package rss2podcast

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/intothevoid/rss2podcast/internal/config"
	"github.com/intothevoid/rss2podcast/internal/io"
	"github.com/intothevoid/rss2podcast/internal/store"
	"github.com/intothevoid/rss2podcast/pkg/audio"
	"github.com/intothevoid/rss2podcast/pkg/fileutil"
	"github.com/intothevoid/rss2podcast/pkg/llm"
	"github.com/intothevoid/rss2podcast/pkg/podcast"
	"github.com/intothevoid/rss2podcast/pkg/rss"
	"github.com/intothevoid/rss2podcast/pkg/tts"
)

type App interface {
	Run() string
}

type rss2podcast struct {
	cfg               *config.Config
	rssParser         *rss.Parser
	store             *store.Store
	ollama            llm.LLM
	podcast           podcast.Podcast
	converter         *tts.Converter
	writer            io.StoreIO
	noConnectionCheck bool
	noParse           bool
	noConvert         bool
	topic             string
}

func NewRSS2Podcast() *rss2podcast {
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

	return &rss2podcast{
		cfg:               cfg,
		rssParser:         rssParser,
		store:             store,
		ollama:            ollama,
		podcast:           podcast,
		converter:         converter,
		writer:            writer,
		noConnectionCheck: noConnectionCheck,
		noParse:           noParse,
		noConvert:         noConvert,
		topic:             "default", //default topic
	}
}

func (r *rss2podcast) SetTopic(topic string) {
	r.topic = topic
}

func (r *rss2podcast) GetTopic() string {
	return r.topic
}

func (r *rss2podcast) Run() (string, error) {
	// Check ollama connection
	if !r.noConnectionCheck {
		log.Println("Checking connection to Ollama...")
		err := r.ollama.CheckConnection()
		if err != nil {
			log.Fatal(err)
			return "", err
		}
	}

	// Clean up old files
	fileutil.CleanupFolder(".", []string{".txt", ".wav", ".mp3", ".json"})

	// Set podcast subject to passed in topic if not default
	if r.topic != "default" {
		r.cfg.Podcast.Subject = r.topic
	}

	// Encode topic for url query
	r.topic = url.QueryEscape(r.topic)

	// Set RSS feed URL to Google News search for topic
	r.cfg.RSS.URL = fmt.Sprintf("https://news.google.com/rss/search?q=%s", r.topic)

	// podcast filename
	// get timestamp as string in format yymmhh_hhmm
	ts := time.Now().Local().Format("2006_01_02_1504")
	podcast_fname_txt := fmt.Sprintf("%s_summary_%s.txt", r.cfg.Podcast.Subject, ts)
	podcast_fname_wav := fmt.Sprintf("%s_summary_%s.wav", r.cfg.Podcast.Subject, ts)
	podcast_fname_mp3 := fmt.Sprintf("%s_summary_%s.mp3", r.cfg.Podcast.Subject, ts)

	tmpStore, err := r.writer.ReadStore()
	if err != nil {
		log.Println("No articles found, starting from scratch")
	} else {
		r.store = tmpStore
	}

	if !r.noParse {
		// Generate podcast introduction
		introduction := "Welcome to the " + r.cfg.Podcast.Subject + " podcast. I'm your host, " +
			r.cfg.Podcast.Podcaster + ". This is an AI generated podcast based on an RSS source feed. " +
			". Let's get started."

		log.Println("Generating podcast introduction...")
		fileutil.AppendStringToFile(podcast_fname_txt, introduction)

		// Parse RSS feed
		items, _ := r.rssParser.ParseURL(r.cfg.RSS.URL)

		// Store top 'maxArticles' articles, default to 10
		for i, item := range items {
			if i >= r.cfg.RSS.MaxArticles {
				break
			}

			rssItem := rss.RSSItem{
				Title:       item.Title,
				Description: item.Description,
				Url:         item.Link,
				HtmlContent: "",
			}
			if !rssItem.IsFiltered(r.cfg.RSS.Filters) {
				r.store.Save(item.GUID, rssItem)
			}
		}
	}

	// Scrape all URLs and populate HTML content
	log.Println("Gathering content from feed websites...")
	r.store.PopulateHtmlContent()

	// Write store to JSON
	r.writer.WriteStore(r.store)

	// Summarize articles
	for _, item := range r.store.GetData() {
		log.Printf("Summarizing article - %s", item.Title)
		summary, err := r.podcast.GenerateSummary(item.Title, item.Description, item.HtmlContent)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		log.Print("Done.")
		fileutil.AppendStringToFile(podcast_fname_txt, summary)
	}

	// Convert podcast text to audio
	if !r.noConvert {
		log.Println("Converting podcast text to audio...")
		fileContent, err := fileutil.ReadFileContent(podcast_fname_txt)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

		// Generate audio file
		r.converter.ConvertToAudio(fileContent, podcast_fname_wav)

		// Convert audio file to mp3
		audio.ConvertWavToMp3(podcast_fname_wav, podcast_fname_mp3)
	}

	return podcast_fname_mp3, nil
}
