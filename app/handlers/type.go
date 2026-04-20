package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type TypeCommand struct {
}

func (c TypeCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	entry, ok := db.Get(key.Bulk)
	if !ok {
		return resp.String("none")
	}

	var response resp.Value
	switch entry.GetType() {
	case database.STRING:
		response = resp.String("string")
	}

	return response
}
