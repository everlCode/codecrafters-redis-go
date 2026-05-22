package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type IncrCommand struct {
}

func (c IncrCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	var digit int
	var entry database.Entry
	entry, ok := db.Get(key.Bulk)
	if !ok {
		digit = 1
		entry = database.Integer(digit)
	}

	value := entry.AsInteger()
	value += 1
	entry.Set(value)
	db.Set(key.Bulk, entry)

	return resp.Integer(value)
}
