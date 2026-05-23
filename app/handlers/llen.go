package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LLenCommand struct {
}

func (c LLenCommand) Execute(args []string, db *database.DB) resp.Value {
	key := args[0]

	value, ok := db.Get(key)
	if !ok || !value.IsArray() {
		return resp.Integer(0)
	}

	lenght := len(value.AsArray())

	return resp.Integer(lenght)
}
