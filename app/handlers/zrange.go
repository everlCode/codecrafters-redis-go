package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ZrangeCommand struct {
}

func (c ZrangeCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 2 {
		return Response(resp.Error("Too few args!"))
	}

	key := args[0]
	startString := args[1]
	start, err := strconv.Atoi(startString)
	if err != nil {
		return Response(resp.Error("Invalid Start"))
	}
	endString := args[2]
	stop, err := strconv.Atoi(endString)
	if err != nil {
		return Response(resp.Error("Invalid End"))
	}

	db := server.GetDB()
	entry, ok := db.Get(key)
	if !ok {
		return Response(resp.EmptyArray())
	}

	zset, ok := entry.AsZset()
	if !ok {
		return Response(resp.EmptyArray())
	}

	lenght := len(zset.Set)
	if start >= lenght || start > stop {
		return Response(resp.EmptyArray())
	}
	
	stop += 1
	if stop > lenght {
		stop = lenght
	}

	var result []string
	var set []database.Zvalue

	if stop == lenght {
		set = zset.Set[start:]
	} else {
		set = zset.Set[start:stop]
	}
	
	for _, v := range set {
		result = append(result, v.Value)
	}

	return Response(resp.ArrayString(result...))
}
