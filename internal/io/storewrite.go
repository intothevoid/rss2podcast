package io

import "github.com/intothevoid/rss2podcast/internal/store"

type StoreWriter interface {
	WriteStore(store *store.Store) error
	ReadStore() (store *store.Store, err error)
}
