package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GetCommand struct {
}

func (c GetCommand) Execute(args []string, server *server.Server) resp.Value {
	db := server.GetDB()
	if len(args) < 1 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := args[0]

	value, ok := db.Get(key)
	if !ok || (value.Expires != 0 && value.Expires < time.Now().UnixMilli()) {
		return resp.Bulk("")
	}
	str := value.AsString()

	return resp.Bulk(str)
}
