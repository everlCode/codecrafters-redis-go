package database

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type DB struct {
	mx      sync.Mutex
	sets    map[string]resp.Value
	waiters map[string][]*Waiter
}

type Waiter struct {
	Chanel  chan resp.Value
	Timeout time.Time
}

func New() *DB {
	return &DB{
		sets:    make(map[string]resp.Value),
		waiters: make(map[string][]*Waiter),
	}
}

func (db *DB) Set(key string, value resp.Value) {
	db.mx.Lock()
	defer db.mx.Unlock()

	if len(db.waiters[key]) > 0 {
		var waiter *Waiter

		for len(db.waiters[key]) > 0 {
			w := db.PopWaiter(key)

			if w.Timeout.IsZero() || time.Now().Before(w.Timeout) {
				waiter = w
				break
			}
		}

		if waiter != nil {
			go func() {
				waiter.Chanel <- value
			}()
			return
		}
	}

	db.sets[key] = value
}

func (db *DB) Get(key string) (resp.Value, bool) {
	db.mx.Lock()
	defer db.mx.Unlock()
	value, ok := db.sets[key]
	if (len(db.sets) == 0) {
		delete(db.sets, key)
	}

	return value, ok
}

func (db *DB) PopWaiter(key string) *Waiter {
	return db.waiters[key][0]
}

func (db *DB) PushWaiter(key string, w *Waiter) {
	db.mx.Lock()
	defer db.mx.Unlock()
	db.waiters[key] = append(db.waiters[key], w)
}
