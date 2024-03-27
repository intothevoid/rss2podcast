package rss

import (
	"github.com/mmcdole/gofeed"
)

type Parser struct {
	fp *gofeed.Parser
}

// NewParser creates a new instance of the Parser struct.
func NewParser() *Parser {
	return &Parser{
		fp: gofeed.NewParser(),
	}
}

// ParseURL parses the RSS feed from the given URL and returns a slice of items and an error, if any.
func (p *Parser) ParseURL(feedURL string) ([]*gofeed.Item, error) {
	feed, err := p.fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}
	return feed.Items, nil
}
