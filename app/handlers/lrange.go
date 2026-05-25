package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
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
	data := entry.AsArray()
	lenght := len(data)

	if start < 0 {
		start = lenght + start
		if start < 0 {
			start = 0
		}
	}

	if end < -1 {
		end = lenght + end + 1
	} else if end < lenght && end > 0 {
		end += 1
	}

	if start > end && end >= 0 {
		return Response(resp.EmptyArray())
	}

	if !ok || !entry.IsArray() {
		return Response(resp.EmptyArray())
	}

	if start >= lenght {
		return Response(resp.EmptyArray())
	}

	if end == -1 || end >= lenght || end == 0 {
		data = data[start:]
	} else {
		data = data[start:end]
	}

	return Response(resp.ArrayString(data))
}