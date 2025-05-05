package store

import (
	"sync"

	"github.com/intothevoid/rss2podcast/pkg/html"
	"github.com/intothevoid/rss2podcast/pkg/rss"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]rss.RSSItem
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]rss.RSSItem),
	}
}

func (s *Store) Save(key string, value rss.RSSItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (rss.RSSItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *Store) GetData() map[string]rss.RSSItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}

// Function to iterate over the data and populate the RSSItem.HtmlContent field
// with the HTML content of the article scraped from the URL
func (s *Store) PopulateHtmlContent() {
	type update struct {
		key   string
		value rss.RSSItem
	}

	updates := make(chan update, len(s.data))
	wg := sync.WaitGroup{}

	// Launch goroutines to scrape content
	for key, item := range s.data {
		wg.Add(1)
		go func(key string, item rss.RSSItem) {
			defer wg.Done()
			item.HtmlContent = html.Scrape(item.Url)
			updates <- update{key: key, value: item}
		}(key, item)
	}

	// Wait for all scraping to complete
	go func() {
		wg.Wait()
		close(updates)
	}()

	// Apply updates to the map
	s.mu.Lock()
	defer s.mu.Unlock()
	for update := range updates {
		s.data[update.key] = update.value
	}
}
