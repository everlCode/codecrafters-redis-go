package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ZscoreCommand struct {
}

func (c ZscoreCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 1 {
		return Response(resp.Error("Too few args!"))
	}

	key := args[0]
	member := args[1]

	db := server.GetDB()
	entry, ok := db.Get(key)
	if !ok {
		return Response(resp.NullBulk())
	}

	zset, ok := entry.AsZset()
	if !ok {
		return Response(resp.NullBulk())
	}
	value, ok := zset.Keys[member]
	if !ok {
		return Response(resp.NullBulk())
	}
	score := strconv.FormatFloat(value.Score, 'f', 10, 64)

	return Response(resp.Bulk(score))
}
