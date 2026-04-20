package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LpushCommand struct {
}

func (c LpushCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	entry, ok := db.Get(key.Bulk)
	if !ok {
		return resp.Array([]string{})
	}
	argss := resp.ParseSlice(args)
	for i := 0; i < len(argss); i++ {
		data := entry.AsArray()
		
		var a = []string{argss[i]}
		entry.Set(append(a, data...))
	}

	arr := entry.AsArray()
	lenght := len(arr)

	db.Set(key.Bulk, entry)

	return resp.Integer(lenght)
}
