package rss

import (
	"github.com/mmcdole/gofeed"
)

type Parser struct {
	fp *gofeed.Parser
}

func NewParser() *Parser {
	return &Parser{
		fp: gofeed.NewParser(),
	}
}

func (p *Parser) ParseURL(feedURL string) ([]*gofeed.Item, error) {
	feed, err := p.fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}
	return feed.Items, nil
}
