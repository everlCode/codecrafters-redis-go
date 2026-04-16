package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LpushCommand struct {
}

func (c LpushCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	arg := args[1:]
	value, ok := db.Sets[key.Bulk]
	if !ok {
		value = resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}
	for i := 0; i < len(arg); i++ {
		var a = []resp.Value{arg[i]}
		value.Array = append(a, value.Array... )
	}
	
	lenght := len(value.Array)

	db.Set(key.Bulk, value)

	return resp.Value{Type: resp.INTEGER, Integer: lenght}
}

func (c LpushCommand) Name() string {
	return LPUSH
}
