package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RpushCommand struct {
}

func (c RpushCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	arg := args[1:]
	entry, ok := db.Get(key.Bulk)
	var value []string
	if !ok {
		entry = database.Array(resp.ParseSlice(arg))
	} else {
		value := entry.AsArray()
		args := resp.ParseSlice(arg)
		entry.Set(append(value[:], args[:]...))
	}

	value = entry.AsArray()
	length := len(value)

	db.Set(key.Bulk, entry)

	return resp.Integer(length)
}
