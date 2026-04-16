package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LLenCommand struct {
}

func (c LLenCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	value, ok := db.Get(key.Bulk)
	if !ok {
		value = resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}
	
	lenght := len(value.Array)

	return resp.Value{Type: resp.INTEGER, Integer: lenght}
}

func (c LLenCommand) Name() string {
	return LLEN
}
