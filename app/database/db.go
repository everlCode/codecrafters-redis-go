package database

import (
	"sync"
	"time"
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
	isMulti bool
	transactionQueue []*CommandQueue
}

type Waiter struct {
	Chanel  chan Entry
	Timeout time.Time
}

func New() *DB {
	return &DB{
		sets:    make(map[string]Entry),
		waiters: make(map[string][]*Waiter),
	}
}

func (db *DB) Set(key string, value Entry) {
	db.mx.Lock()
	defer db.mx.Unlock()

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

func (db *DB) SetMulti(value bool) {
	db.isMulti = value
}

func (db *DB) IsTransaction() bool {
	return db.isMulti
}

func (db *DB) PushTransactionQueue(command CommandQueue)  {
	db.transactionQueue = append(db.transactionQueue, &command)
}

func (db *DB) PopTransactionQueue(command CommandQueue) *CommandQueue {
	if len(db.transactionQueue) > 0 {
		return db.transactionQueue[0]
	}

	return nil
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

func Array(data []string) Entry {
	return Entry{value: data}
}

func Integer(data int) Entry {
	return Entry{value: data}
}

func String(data string) Entry {
	return Entry{value: data}
}
