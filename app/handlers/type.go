package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type TypeCommand struct {
}

func (c TypeCommand) Execute(args []string, server *server.Server, client *clients.Client) resp.Value {
	db := server.GetDB()
	key := args[0]

	entry, ok := db.Get(key)
	if !ok {
		return resp.SimpleString("none")
	}

	var response string
	switch entry.GetType() {
	case database.STRING:
		response = "string"
	case database.STREAM:
		response = "stream"
	default:
		response = "undefined"
	}

	return resp.SimpleString(response)
}
