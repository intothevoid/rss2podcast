package io

import "github.com/intothevoid/rss2podcast/internal/store"

type StoreIO interface {
	WriteStore(store *store.Store) error
	ReadStore() (store *store.Store, err error)
}
