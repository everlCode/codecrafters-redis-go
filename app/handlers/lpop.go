package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type LPopCommand struct {
}

func (c LPopCommand) Execute(args []string, server *server.Server) resp.Value {
	db := server.GetDB()
	key := args[0]

	var count int = 1
	if len(args) > 1 {
		v := args[1]
		c, err := strconv.Atoi(v)
		count = c
		if err != nil {
			return resp.Value{Type: resp.ERROR}
		}
	}

	entry, ok := db.Get(key)
	var response resp.Value
	if !ok || !entry.IsArray() {
		return resp.Value{
			Type: resp.BULK,
		}
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
		response = resp.ArrayString(v)
	}
	entry.Set(value[count:])

	db.Set(key, entry)

	return response
}
