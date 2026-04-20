package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LPopCommand struct {
}

func (c LPopCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	var count int = 1
	if len(args) > 1 {
		v := args[1]
		if v.Type != resp.BULK {
			return resp.Value{Type: resp.ERROR}
		}
		c, err := strconv.Atoi(v.Bulk)
		count = c
		if err != nil {
			return resp.Value{Type: resp.ERROR}
		}
	}

	entry, ok := db.Get(key.Bulk)
	var response resp.Value
	if !ok || !entry.IsArray() {
		return resp.Value{
			Type: resp.BULK,
		}
	}

	value, _ := entry.AsArray()
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
		response = resp.Array(v)
	}
	entry.Set(value[count:])

	db.Set(key.Bulk, entry)

	return response
}
