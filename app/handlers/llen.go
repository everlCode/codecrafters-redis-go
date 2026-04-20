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
	var response resp.Value
	if !ok || value.IsArray() {
		response = resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	lenght := len(response.Array)

	return resp.Integer(lenght)
}
