package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type XaddCommand struct {
}

func (c XaddCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := args[0]

	entry, ok := db.Get(key.Bulk)
	if !ok || (entry.Expires != 0 && entry.Expires < time.Now().UnixMilli()) {
		return resp.Value{Type: resp.BULK, Bulk: ""}
	}

	return resp.Value{}
}
