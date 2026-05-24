package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type EchoCommand struct {
}

func (c EchoCommand) Execute(args []string, server *server.Server) resp.Value {
	if len(args) == 0 {
		return resp.Bulk("PONG")
	}

	return resp.Bulk(args[0])
}
