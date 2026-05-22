package handlers

import (
	"strconv"

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

	entry, ok := db.Get(key.Bulk)
	if !ok {
		return resp.Integer(0)
	}

	value := entry.AsString()
	digit, err := strconv.Atoi(value)
	if err != nil {
		return resp.Integer(0)
	}
	digit += 1
	entry.Set(digit)
	db.Set(key.Bulk, entry)

	return resp.Integer(digit)
}
