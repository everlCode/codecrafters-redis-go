package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LpushCommand struct {
}

func (c LpushCommand) Execute(args []string, db *database.DB, client *clients.Client) resp.Value {
	key := args[0]
	
	entry, ok := db.Get(key)
	if !ok {
		entry = database.Array([]string{})
	}

	argss := args[1:]
	for i := range argss {
		data := entry.AsArray()
		
		var a = []string{argss[i]}
		entry.Set(append(a, data...))
	}
	

	arr := entry.AsArray()
	lenght := len(arr)

	db.Set(key, entry)

	return resp.Integer(lenght)
}
