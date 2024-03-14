package io

import (
	"encoding/json"
	"os"

	"github.com/intothevoid/rss2podcast/internal/store"
	"github.com/intothevoid/rss2podcast/pkg/rss"
)

type writer struct {
	store *store.Store
}

func NewJsonWriter(store *store.Store) StoreIO {
	return &writer{
		store: store,
	}
}

func (w *writer) WriteStore(store *store.Store) error {
	data := w.store.GetData()
	file, err := os.Create("articles.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for key, rssItem := range data {
		// Create a map with the uid as the key and the rssItem as the value
		itemMap := map[string]rss.RSSItem{key: rssItem}
		// Encode the map to JSON and write it to the file
		if err := encoder.Encode(itemMap); err != nil {
			return err
		}
	}

	return nil
}

func (w *writer) ReadStore() (*store.Store, error) {
	file, err := os.Open("articles.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	store := store.NewStore()

	for decoder.More() {
		var data map[string]rss.RSSItem
		if err := decoder.Decode(&data); err != nil {
			return nil, err
		}

		for key, rssItem := range data {
			store.Save(key, rssItem)
		}
	}

	return store, nil
}
