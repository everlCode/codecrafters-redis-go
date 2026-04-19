package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type EchoCommand struct {
}

func (c EchoCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.BULK, Bulk: "PONG"}
	}

	return resp.Value{Type: resp.BULK, Bulk: args[0].Bulk}
}
