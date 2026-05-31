package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ZaddCommand struct {
}

func (c ZaddCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 2 {
		return Response(resp.Error("Too few args!"))
	}

	key := args[0]
	idString := args[1]
	id, err := strconv.ParseFloat(idString, 64)
	if err != nil {
		return Response(resp.Error("Invalid Id"))
	}
	value := args[2]

	db := server.GetDB()
	entry, ok := db.Get(key)
	if !ok {
		zset := database.NewZset()
		entry = *database.NewEntry()
		entry.Set(zset)
	}

	zset, ok := entry.AsZset()
	if !ok {
		return Response(resp.Error("Key" + key + " doenst sorted set!"))
	}
	isNewValue := zset.Add(id, value)
	entry.Set(zset)
	db.Set(key, entry)

	return Response(resp.Integer(isNewValue))
}
