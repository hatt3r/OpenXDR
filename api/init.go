package api

import "openxdr/internal/store"

var DB *store.Store

func Init(db *store.Store) {
	DB = db
}
