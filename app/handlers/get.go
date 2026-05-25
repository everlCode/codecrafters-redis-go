package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GetCommand struct {
}

func (c GetCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	db := server.GetDB()
	if len(args) < 1 {
		return Response(resp.Error("ERR to few args"))
	}

	key := args[0]

	value, ok := db.Get(key)
	if !ok || (value.Expires != 0 && value.Expires < time.Now().UnixMilli()) {
		return Response(resp.NullBulk())
	}
	str := value.AsString()

	return Response(resp.Bulk(str))
}
