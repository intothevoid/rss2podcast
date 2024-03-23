package main

import (
	rss2podcast "github.com/intothevoid/rss2podcast/internal/app"
)

func main() {
	// Create application
	app := rss2podcast.NewRSS2Podcast()

	// Run application
	app.Run()
}
