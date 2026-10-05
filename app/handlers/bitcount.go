package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type BitCountCommand struct {
}

func (c BitCountCommand) Execute(
	args []string,
	server *server.Server,
	client *clients.Client,
) CommandResponse {

	db := server.GetDB()
	if len(args) < 1 {
		return Response(resp.Integer(0))
	}
	key := args[0]
	var start, end int

	if len(args) > 2 {
		s, err := strconv.Atoi(args[1])
		if err == nil {
			start = s
		}
		e, err := strconv.Atoi(args[2])
		if err == nil {
			end = e
		}
	}
	
	var entry database.Entry
	entry, ok := db.Get(key);
	
	if !ok {
		entry = database.String("")
	}
	arr := []byte(entry.AsString())
	if len(args) > 2 {
		arr =  helpers.GetDataByStartEnd(arr, start, end)
	}
	
	count := helpers.CountTrueBit(arr)

	response := Response(resp.Integer(count))

	return response
}
