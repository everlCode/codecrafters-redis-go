package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type GetCommand struct {
}

func (c GetCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := args[0]

	value, ok := db.Get(key.Bulk)
	if !ok || (value.Expires != 0 && value.Expires < time.Now().UnixMilli()) {
		return resp.Value{Type: resp.BULK, Bulk: ""}
	}
	str := value.AsString()

	return resp.SimpleString(str)
}
