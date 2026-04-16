package db

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type DB struct {
	sets map[string]resp.Value
}

func New() *DB {
	return &DB{
		sets: make(map[string]resp.Value),
	}
}

func (db *DB) Set(key string, value resp.Value) {
	db.sets[key] = value
}

func (db *DB) Get(key string) (resp.Value, bool) {
	value, ok := db.sets[key]

	return value, ok
}
