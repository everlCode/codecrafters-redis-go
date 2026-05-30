package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type KeysCommand struct {
}

func (c KeysCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	db := server.GetDB()
	set := db.GetAll()

	var keys []string
	for key := range set {
		keys = append(keys, key)
	}

	return Response(resp.ArrayString(keys...))
}
