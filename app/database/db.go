package database

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	UNKNOWN = 0
	STRING  = 1
	ARRAY   = 2
	STREAM  = 2
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
	id   string
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

func (v Entry) AsString() string {
	a, ok := v.value.(string)
	if !ok {
		panic(fmt.Sprintf("value is not string: %T", v.value))
	}
	return a
}

func (v Entry) AsArray() []string {
	a, _ := v.value.([]string)
	return a
}

func (v Entry) AsStream() (Stream, bool) {
	a, ok := v.value.(Stream)
	return a, ok
}

func (s Stream) GetLastId() string {
	return s.id
}

func (s Stream) GetLastIdParts() (int, int) {
	parts := strings.Split(s.id, "-")
	if len(parts) < 2 {
		return 0, 0
	}

	miliseconds, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0
	}

	return miliseconds, id
}

func CreateStream(id string, data map[string]string) Stream {
	stream := NewStream()
	stream.Add(id, data)

	return stream
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
	case Stream:
		return STREAM
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
	stream.id = id
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
