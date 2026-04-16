package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LPopCommand struct {
}

func (c LPopCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	var count int = 1
	if (len(args) > 1) {
		v := args[1]
		if v.Type != resp.BULK {
			return resp.Value{Type: resp.ERROR} 
		}
		c, err := strconv.Atoi(v.Bulk)
		count = c
		if err != nil  {
			return resp.Value{Type: resp.ERROR}
		}
	}
	
	value, ok := db.Get(key.Bulk)
	if !ok {
		value = resp.Value{
			Type: resp.BULK,
		}

		return value
	}

	lenght := len(value.Array)
	if (count > lenght) {
		count = lenght
	}

	var response resp.Value
	if (count == 1) {
		v := value.Array[0]
		response.Type = resp.BULK
		response.Bulk = v.Bulk
	} else {
		v := value.Array[:count]
		response.Type = resp.ARRAY
		response.Array = v
	}
	value.Array = value.Array[count:]
	
	db.Set(key.Bulk, value)

	return response
}

func (c LPopCommand) Name() string {
	return LPOP
}
