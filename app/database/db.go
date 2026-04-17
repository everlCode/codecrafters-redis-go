package database

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type DB struct {
	sets    map[string]resp.Value
	waiters map[string][]Waiter
}

type Waiter struct {
	Chanel chan resp.Value
}

func New() *DB {
	return &DB{
		sets:    make(map[string]resp.Value),
		waiters: make(map[string][]Waiter),
	}
}

func (db *DB) Set(key string, value resp.Value) {
	if len(db.waiters[key]) > 0 {
		waiter := db.PopWaiter(key)
		waiter.Chanel <- value
		return
	}

	db.sets[key] = value
}

func (db *DB) Get(key string) (resp.Value, bool) {
	value, ok := db.sets[key]

	return value, ok
}

func (db *DB) PopWaiter(key string) Waiter {
	return db.waiters[key][0]
}

func (db *DB) PushWaiter(key string, w Waiter) {
	db.waiters[key] = append(db.waiters[key], w)
}
