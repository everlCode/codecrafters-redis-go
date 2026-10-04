package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type SetBitCommand struct {
}

func (c SetBitCommand) Execute(
	args []string,
	server *server.Server,
	client *clients.Client,
) CommandResponse {

	db := server.GetDB()
	if len(args) < 3 {
		return Response(resp.Error("ERR to few args"))
	}

	key := args[0]
	var offset int
	offset, err := strconv.Atoi(args[1]);
	if err != nil {
		return Response(resp.Error("Invalid offset!"))
	}
	value, err := strconv.Atoi(args[2])
	if err != nil || value < 0 || value > 1 {
		return Response(resp.Error("Invalid value!"))
	}
	var entry database.Entry
	entry, ok := db.Get(key);
	
	if !ok {
		entry = database.String("")
	}
	
	bytes, originaValue, err := helpers.SetBit([]byte(entry.AsString()), offset, value != 0)
	entry.Set(string(bytes))

	db.Set(key, entry)

	response := Response(resp.Integer(originaValue))
	response.SetPropagationCommand(SETBIT, args)

	return response
}
