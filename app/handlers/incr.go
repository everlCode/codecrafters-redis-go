package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type IncrCommand struct {
}

func (c IncrCommand) Execute(args []string, db *database.DB, client *clients.Client) resp.Value {
	key := args[0]

	var entry database.Entry
	entry, ok := db.Get(key)
	if !ok {
		entry = database.String("0")
	}

	value := entry.AsString()
	v, err := strconv.Atoi(value)
	if err != nil {
		return resp.Error("ERR value is not an integer or out of range")
	}
	v += 1
	value = strconv.Itoa(v)

	entry.Set(value)
	db.Set(key, entry)

	return resp.Integer(v)
}
