package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type EchoCommand struct {
}

func (c EchoCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.BULK, Bulk: "PONG"}
	}

	return resp.Value{Type: resp.BULK, Bulk: args[0].Bulk}
}

func (c EchoCommand) Name() string {
	return ECHO
}
