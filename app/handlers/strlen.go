package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type StrLenCommand struct {
}

func (c StrLenCommand) Execute(
	args []string,
	server *server.Server,
	client *clients.Client,
) CommandResponse {

	db := server.GetDB()
	if len(args) < 1 {
		return Response(resp.Error("ERR to few args"))
	}

	key := args[0]
	
	var entry database.Entry
	entry, ok := db.Get(key);
	if !ok {
		return Response(resp.Integer(0))
	}

	response := Response(resp.Integer(len(entry.AsString())))

	return response
}
