package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	PING = "PING"
	SET  = "SET"
	GET  = "GET"
	ECHO = "ECHO"
)

type Command interface {
	Execute([]resp.Value, *db.DB) resp.Value
	Name() string
}

func echo(args []resp.Value, db *db.DB) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.BULK, Bulk: "PONG"}
	}

	return resp.Value{Type: resp.BULK, Bulk: args[0].Bulk}
}
