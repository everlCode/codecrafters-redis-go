package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type LLenCommand struct {
}

func (c LLenCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	db := server.GetDB()
	key := args[0]

	value, ok := db.Get(key)
	if !ok || !value.IsArray() {
		return Response(resp.Integer(0))
	}

	lenght := len(value.AsArray())

	return Response(resp.Integer(lenght))
}
