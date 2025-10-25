package main

type DB struct {
	sets map[string]Value
}

func NewDB() *DB {
	return &DB{
		sets: make(map[string]Value),
	}
}