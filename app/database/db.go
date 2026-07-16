package database

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
)

const (
	UNKNOWN = 0
	STRING  = 1
	ARRAY   = 2
	STREAM  = 3
)

type DB struct {
	mx      sync.Mutex
	sets    map[string]Entry
	waiters map[string][]*Waiter
	WatchedKeys map[string]uint64
}

type Waiter struct {
	Chanel  chan Entry
	Timeout time.Time
}

func New(config *config.Config) *DB {
	rd := rdb.New(config)
	var data []rdb.Value
	if rd != nil {
		data = rd.Read()
	}
	
	db :=  &DB{
		sets:    make(map[string]Entry),
		waiters: make(map[string][]*Waiter),
		WatchedKeys: make(map[string]uint64),
	}
	db.ImportRDB(data)

	return db
}

func (db *DB) Set(key string, value Entry) {
	db.mx.Lock()
	defer db.mx.Unlock()
	db.WatchedKeys[key]++

	if len(db.waiters[key]) > 0 {
		var waiter *Waiter

		waiter = db.PopWaiter(key)
	
		if waiter != nil {
			go func() {
				waiter.Chanel <- value
			}()
			return
		}
	}

	db.sets[key] = value
}

func (db *DB) Get(key string) (Entry, bool) {
	db.mx.Lock()
	defer db.mx.Unlock()
	value, ok := db.sets[key]
	if len(db.sets) == 0 {
		delete(db.sets, key)
	}

	return value, ok
}

func (db *DB) GetAll() map[string]Entry {
	return db.sets
}

func (db *DB) IsWaiterForKeyExist(key string) bool {
	_, ok := db.waiters[key]

	return ok
}

func (db *DB) PopWaiter(key string) *Waiter {
	for  _,w := range db.waiters[key] {
		if w.Timeout.IsZero() || time.Now().Before(w.Timeout) {
			return w
		}
	}

	return nil
}

func (db *DB) PushWaiter(key string, time time.Time) chan Entry {
	db.mx.Lock()
	defer db.mx.Unlock()

	ch := make(chan Entry)
	w := Waiter{
		Timeout: time,
		Chanel: ch,
	}
	db.waiters[key] = append(db.waiters[key], &w)

	return ch
}

func (db *DB) GetVersion(key string) uint64 {
	return db.WatchedKeys[key]
}

func (db *DB) ImportRDB(data []rdb.Value) {
	for _, v := range data {
		var entry Entry
		switch v.Type {
		case rdb.STRING:
			val, _ := v.Value.(string)
			entry = String(val)
			entry.Expires = int64(v.Expire)
		
		}
		db.Set(v.Key, entry)	
	}
}

func Array(data []string) Entry {
	return Entry{value: data}
}

func Integer(data int) Entry {
	return Entry{value: data}
}

func String(data string) Entry {
	return Entry{
		value: data,
	}
}
