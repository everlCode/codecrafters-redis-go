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
