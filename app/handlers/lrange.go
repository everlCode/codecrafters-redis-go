package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type LRangeCommand struct {
}

func (c LRangeCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	db := server.GetDB()
	if len(args) < 3 {
		return Response(resp.Error("ERR to few args"))
	}
	key := args[0]
	startValue := args[1]
	endValue := args[2]

	start, err := strconv.Atoi(startValue)
	if err != nil {
		return Response(resp.Error(err.Error()))
	}
	end, err := strconv.Atoi(endValue)
	if err != nil {
		return Response(resp.Error(err.Error()))
	}

	entry, ok := db.Get(key)
	if !ok || !entry.IsArray() {
		return Response(resp.EmptyArray())
	}
	data := entry.AsArray()

	data = helpers.GetDataByStartEnd(data, start, end)
	if len(data) == 0 {
		return Response(resp.EmptyArray())
	}

	return Response(resp.ArrayString(data...))
}