package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LPopCommand struct {
}

func (c LPopCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	value, ok := db.Get(key.Bulk)
	if !ok {
		value = resp.Value{
			Type: resp.BULK,
		}

		return value
	}
	v := value.Array[0]
	value.Array = value.Array[1:]
	db.Set(key.Bulk, value)

	return resp.Value{Type: resp.BULK, Bulk: v.Bulk}
}

func (c LPopCommand) Name() string {
	return LPOP
}
