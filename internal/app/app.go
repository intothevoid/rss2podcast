package rss2podcast

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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
	noSummary         bool
	topic             string
	topics            []string
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

	coquiConfig := &tts.ConverterConfig{
		URL: cfg.TTS.Coqui.URL,
	}

	kokoroConfig := &tts.ConverterConfig{
		URL:    cfg.TTS.Kokoro.URL,
		Voice:  cfg.TTS.Kokoro.Voice,
		Speed:  cfg.TTS.Kokoro.Speed,
		Format: cfg.TTS.Kokoro.Format,
	}

	mlxConfig := &tts.ConverterConfig{
		URL:    cfg.TTS.MLX.URL,
		Voice:  cfg.TTS.MLX.Voice,
		Speed:  cfg.TTS.MLX.Speed,
		Format: cfg.TTS.MLX.Format,
	}

	converter := tts.NewConverter(cfg.TTS.Engine, coquiConfig, kokoroConfig, mlxConfig)
	writer := io.NewJsonWriter(store)

	// Check command line arguments
	noParse := false
	noConvert := false
	noConnectionCheck := false
	noSummary := false
	var topics []string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--no-parse":
			noParse = true
		case "--no-convert":
			noConvert = true
		case "--no-connection-check":
			noConnectionCheck = true
		case "--no-summary":
			noSummary = true
		case "--topics":
			if i+1 < len(os.Args) {
				topics = strings.Split(os.Args[i+1], ",")
				for i, topic := range topics {
					topics[i] = strings.TrimSpace(topic)
				}
				i++ // Skip the next argument since we've processed it
			}
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
		noSummary:         noSummary,
		topic:             "default", //default topic
		topics:            topics,
	}
}

func (r *rss2podcast) SetTopic(topic string) {
	r.topic = topic
}

func (r *rss2podcast) GetTopic() string {
	return r.topic
}

func (r *rss2podcast) Run() (string, error) {
	// Move old mp3 files to an archive folder
	fileutil.MoveFilesToArchiveFolder(".", []string{".mp3"})

	if len(r.topics) > 0 {
		var lastPodcast string
		for _, topic := range r.topics {
			r.SetTopic(topic)
			podcast, err := r.runSingleTopic()
			if err != nil {
				return "", err
			}
			lastPodcast = podcast
		}
		return lastPodcast, nil
	}
	return r.runSingleTopic()
}

func (r *rss2podcast) runSingleTopic() (string, error) {
	// Check ollama connection
	if !r.noConnectionCheck {
		log.Println("Checking connection to Ollama...")
		err := r.ollama.CheckConnection()
		if err != nil {
			log.Fatal(err)
			return "", err
		}
	}

	// Clean up old files, .txt, .wav, .mp3, .json only if not --noParse
	if !r.noParse {
		fileutil.CleanupFolder(".", []string{".txt", ".wav", ".json"})
	} else {
		// Since we're not parsing, we re-use the old json and txt files
		fileutil.CleanupFolder(".", []string{".wav"})
	}

	// Set podcast subject to passed in topic if not default
	if strings.TrimSpace(strings.ToLower(r.topic)) != "default" {
		r.cfg.Podcast.Subject = r.topic
	}

	// Encode topic for url query
	r.topic = url.QueryEscape(r.topic)

	// Set RSS feed URL to Google News search for topic
	r.cfg.RSS.URL = fmt.Sprintf("https://flipboard.com/topic/%s.rss", r.cfg.Podcast.Subject)
	// r.cfg.RSS.URL = fmt.Sprintf("https://news.google.com/rss/search?q=%s", r.topic)

	// Find existing txt file if --no-parse is set
	var existingTxtFile string
	if r.noParse {
		files, err := filepath.Glob(fmt.Sprintf("%s_summary_*.txt", r.cfg.Podcast.Subject))
		if err == nil && len(files) > 0 {
			existingTxtFile = files[0]
		}
	}

	// podcast filename
	// get timestamp as string in format yymmhh_hhmm
	ts := time.Now().Local().Format("2006_01_02_1504")
	if existingTxtFile != "" {
		// Extract timestamp from existing file
		re := regexp.MustCompile(`_(\d{4}_\d{2}_\d{2}_\d{4})\.txt$`)
		if matches := re.FindStringSubmatch(existingTxtFile); len(matches) > 1 {
			ts = matches[1]
		}
	}
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
		introduction := "Hello! Welcome to the " + r.cfg.Podcast.Subject + " podcast. I'm your host, " +
			r.cfg.Podcast.Podcaster + ". This is an AI podcast generated from information on the internet. " +
			"Thanks for tuning in."

		log.Println("Generating podcast introduction...")
		log.Printf("Writing introduction to file: %s", podcast_fname_txt)
		fileutil.FlushStringToFile(podcast_fname_txt, introduction)
		// Parse RSS feed
		items, _ := r.rssParser.ParseURL(r.cfg.RSS.URL)

		// Sort by publication date to keep the most recent articles first
		sort.Slice(items, func(i, j int) bool {
			return items[i].PublishedParsed.After(*items[j].PublishedParsed)
		})

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

		// Scrape all URLs and populate HTML content
		log.Println("Gathering content from feed websites...")
		r.store.PopulateHtmlContent()

		// Write store to JSON
		r.writer.WriteStore(r.store)
	}

	if !r.noSummary {
		// Buffer to store summaries in memory
		summaryBuffer := make(map[string]string)

		// Summarize articles
		for _, item := range r.store.GetData() {
			log.Printf("Summarizing article - %s", item.Title)
			summary, err := r.podcast.GenerateSummary(item.Title, item.Description, item.HtmlContent)
			if err != nil {
				log.Fatal(err)
				return "", err
			}
			log.Print("Done.")
			summaryBuffer[item.Title] = summary
		}
		log.Printf("Writing summaries to file: %s", podcast_fname_txt)
		fileutil.FlushMapToFile(podcast_fname_txt, summaryBuffer)
	}

	// Convert podcast text to audio
	if !r.noConvert {
		log.Println("Converting podcast text to audio...")
		// Check if the text file exists before attempting to read it
		if _, err := os.Stat(podcast_fname_txt); os.IsNotExist(err) {
			log.Printf("Text file %s does not exist. Skipping audio conversion.", podcast_fname_txt)
			return "", fmt.Errorf("text file %s does not exist", podcast_fname_txt)
		}

		fileContent, err := fileutil.ReadFileContent(podcast_fname_txt)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

		// Generate audio file
		r.converter.ConvertToAudio(fileContent, podcast_fname_wav)

		// Convert audio file to mp3
		audio.ConvertWavToMp3(podcast_fname_wav, podcast_fname_mp3)
	} else {
		log.Println("Skipping audio conversion. --noConvert was passed.")
	}

	return podcast_fname_mp3, nil
}
