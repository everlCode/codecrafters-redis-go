package database

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type DB struct {
	mx sync.Mutex
	sets    map[string]resp.Value
	waiters map[string][]Waiter
}

type Waiter struct {
	Chanel chan resp.Value
	Timeout time.Time
}

func New() *DB {
	return &DB{
		sets:    make(map[string]resp.Value),
		waiters: make(map[string][]Waiter),
	}
}

func (db *DB) Set(key string, value resp.Value) {
	db.mx.Lock()
	
	if len(db.waiters[key]) > 0 {
		for i := 0; i < len(db.waiters[key]); i++ {
			waiter := db.PopWaiter(key)
			if waiter.Timeout.IsZero() || !time.Now().After(waiter.Timeout) {
				db.mx.Unlock()
				waiter.Chanel <- value
				break
			} else {
				db.mx.Unlock()
			}
		}
		
		return
	}
	db.mx.Unlock()
	db.sets[key] = value
}

func (db *DB) Get(key string) (resp.Value, bool) {
	db.mx.Lock()
	defer db.mx.Unlock()
	value, ok := db.sets[key]

	return value, ok
}

func (db *DB) PopWaiter(key string) Waiter {
	return db.waiters[key][0]
}

func (db *DB) PushWaiter(key string, w Waiter) {
	db.mx.Lock()
	defer db.mx.Unlock()
	db.waiters[key] = append(db.waiters[key], w)
}
