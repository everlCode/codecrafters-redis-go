package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ZcardCommand struct {
}

func (c ZcardCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 1 {
		return Response(resp.Error("Too few args!"))
	}

	key := args[0]

	db := server.GetDB()
	entry, ok := db.Get(key)
	if !ok {
		return Response(resp.Integer(0))
	}

	zset, ok := entry.AsZset()
	if !ok {
		return Response(resp.Integer(0))
	}

	return Response(resp.Integer(len(zset.Set)))
}
