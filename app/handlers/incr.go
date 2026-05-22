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

	var entry database.Entry
	entry, ok := db.Get(key.Bulk)
	if !ok {
		entry = database.String("0")
	}

	value := entry.AsString()
	v, _ := strconv.Atoi(value)
	v += 1
	value = strconv.Itoa(v)

	entry.Set(value)
	db.Set(key.Bulk, entry)

	return resp.Integer(v)
}
