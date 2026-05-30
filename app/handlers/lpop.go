package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type LPopCommand struct {
}

func (c LPopCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	db := server.GetDB()
	key := args[0]

	var count int = 1
	if len(args) > 1 {
		v := args[1]
		c, err := strconv.Atoi(v)
		count = c
		if err != nil {
			return Response(resp.Error(""))
		}
	}

	entry, ok := db.Get(key)
	var response resp.Value
	if !ok || !entry.IsArray() {
		return Response(resp.Bulk(""))
	}

	value := entry.AsArray()
	lenght := len(value)
	if count > lenght {
		count = lenght
	}

	if count == 1 {
		v := value[0]
		response.Type = resp.BULK
		response.Bulk = v
	} else {
		v := value[:count]
		response = resp.ArrayString(v...)
	}
	entry.Set(value[count:])

	db.Set(key, entry)

	return Response(response)
}
