package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ZrankCommand struct {
}

func (c ZrankCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 1 {
		return Response(resp.Error("Too few args!"))
	}

	key := args[0]
	value := args[1]

	db := server.GetDB()
	entry, ok := db.Get(key)
	if !ok {
		return Response(resp.NullBulk())
	}

	zset, ok := entry.AsZset()
	if !ok {
		return Response(resp.Error("Key" + key + " doenst sorted set!"))
	}
	index, ok := zset.GetIndex(value)
	if !ok {
		return Response(resp.NullBulk())
	}

	return Response(resp.Integer(index))
}
