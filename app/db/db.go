package db

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type DB struct {
	Sets map[string]resp.Value
}

func New() *DB {
	return &DB{
		Sets: make(map[string]resp.Value),
	}
}

func (db *DB) Set(key string, value resp.Value) {
	db.Sets[key] = value
}

func (db *DB) Get(key string) (resp.Value, bool) {
	value, ok := db.Sets[key]

	return value, ok
}
