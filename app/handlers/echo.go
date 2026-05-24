package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type EchoCommand struct {
}

func (c EchoCommand) Execute(args []string, server *server.Server) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.BULK, Bulk: "PONG"}
	}

	return resp.Value{Type: resp.BULK, Bulk: args[0]}
}
