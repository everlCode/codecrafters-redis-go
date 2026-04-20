package database

import (
	"sync"
	"time"
)

const (
	UNKNOWN = 0
	STRING  = 1
	ARRAY   = 2
	STREAM   = 2
)

type DB struct {
	mx      sync.Mutex
	sets    map[string]Entry
	waiters map[string][]*Waiter
}

type Entry struct {
	value   any
	Expires int64
}

type Stream struct {
	data map[string]map[string]string
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

func (db *DB) Get(key string) (Entry, bool) {
	db.mx.Lock()
	defer db.mx.Unlock()
	value, ok := db.sets[key]
	if len(db.sets) == 0 {
		delete(db.sets, key)
	}

	return value, ok
}

func (v Entry) AsString() (string, bool) {
	a, ok := v.value.(string)
	return a, ok
}

func (v Entry) AsArray() ([]string, bool) {
	a, ok := v.value.([]string)
	return a, ok
}

func (v Entry) AsStream() (Stream, bool) {
	a, ok := v.value.(Stream)
	return a, ok
}

func CreateStream(id string, data map[string]string) Entry {
	entry := Entry{}
	stream := NewStream()
	stream.Add(id, data)
	entry.Set(stream)

	return entry
}

func NewStream() Stream {
    return Stream{
        data: make(map[string]map[string]string),
    }
}

func (v Entry) GetType() int {
	switch v.value.(type) {
	case string:
		return STRING
	case []string:
		return ARRAY
	default:
		return UNKNOWN
	}
}

func (e Entry) IsArray() bool {
	return e.GetType() == ARRAY
}

func (e Entry) IsExpired() bool {
	return e.Expires != 0 && e.Expires < time.Now().UnixMilli()
}

func Array(data []string) Entry {
	return Entry{value: data}
}

func (v *Entry) Set(a any) {
	v.value = a
}

func (stream *Stream) Add(id string, data map[string]string) {
	stream.data[id] = data
}

func (db *DB) PopWaiter(key string) *Waiter {
	waiter := db.waiters[key][0]
	if len(db.waiters[key]) > 0 {
		db.waiters[key] = db.waiters[key][1:]
	}

	return waiter
}

func (db *DB) PushWaiter(key string, w *Waiter) {
	db.mx.Lock()
	defer db.mx.Unlock()
	db.waiters[key] = append(db.waiters[key], w)
}
