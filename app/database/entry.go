package database

import (
	"fmt"
	"time"
)

type Entry struct {
	value   any
	Expires int64
}

func (v Entry) AsString() string {
	a, ok := v.value.(string)
	if !ok {
		panic(fmt.Sprintf("value is not string: %T", v.value))
	}
	return a
}

func (v Entry) AsInteger() int {
	a, ok := v.value.(int)
	if !ok {
		panic(fmt.Sprintf("value is not int: %T", v.value))
	}
	return a
}

func (v Entry) AsArray() []string {
	a, _ := v.value.([]string)
	return a
}

func (v Entry) AsStream() (*Stream, bool) {
	a, ok := v.value.(*Stream)
	return a, ok
}

func (v Entry) GetType() int {
	switch v.value.(type) {
	case string:
		return STRING
	case []string:
		return ARRAY
	case *Stream:
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

func (v *Entry) Set(a any) {
	v.value = a
}
