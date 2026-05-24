package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type PingCommand struct {
}

func (c PingCommand) Execute(args []string, server *server.Server) resp.Value {
	return resp.Value{Type: resp.STRING, String: "PONG"}
}
