package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type GetCommand struct {
}

func (c GetCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := args[0]

	value, ok := db.Sets[key.Bulk]
	if !ok || (value.Expires != 0 && value.Expires < time.Now().UnixMilli()) {
		return resp.Value{Type: resp.BULK, Bulk: ""}
	}

	return value
}

func (c GetCommand) Name() string {
	return GET
}
