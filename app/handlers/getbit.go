package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GetBitCommand struct {
}

func (c GetBitCommand) Execute(
	args []string,
	server *server.Server,
	client *clients.Client,
) CommandResponse {

	db := server.GetDB()
	if len(args) < 2 {
		return Response(resp.Error("ERR to few args"))
	}

	key := args[0]
	var offset int
	offset, err := strconv.Atoi(args[1]);
	if err != nil {
		return Response(resp.Error("Invalid offset!"))
	}
	
	var entry database.Entry
	entry, ok := db.Get(key);
	if !ok {
		return Response(resp.Integer(0))
	}
	
	originaValue := helpers.GetBit([]byte(entry.AsString()), offset)

	response := Response(resp.Integer(originaValue))

	return response
}
