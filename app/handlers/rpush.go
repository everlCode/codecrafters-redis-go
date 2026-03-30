package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RpushCommand struct {
}

func (c RpushCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	arg := args[1]
	value, ok := db.Sets[key.Bulk]
	if !ok {
		value = resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{arg},
		}
	} else {
		value.Array = append(value.Array, arg)
	}
	
	len := len(value.Array)

	db.Sets[key.Bulk] = value

	return resp.Value{Type: resp.INTEGER, Integer: len}
}

func (c RpushCommand) Name() string {
	return RPUSH
}
