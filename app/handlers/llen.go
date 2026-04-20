package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LLenCommand struct {
}

func (c LLenCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	value, ok := db.Get(key.Bulk)
	if !ok || !value.IsArray() {
		return resp.Integer(0)
	}

	lenght := len(value.AsArray())

	return resp.Integer(lenght)
}
