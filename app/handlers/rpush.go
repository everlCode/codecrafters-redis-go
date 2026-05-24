package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RpushCommand struct {
}

func (c RpushCommand) Execute(args []string, db *database.DB) resp.Value {
	key := args[0]

	args = args[1:]
	entry, ok := db.Get(key)
	var value []string
	if !ok {
		entry = database.Array(args)
	} else {
		value := entry.AsArray()
		entry.Set(append(value[:], args[:]...))
	}

	value = entry.AsArray()
	length := len(value)

	db.Set(key, entry)

	return resp.Integer(length)
}
