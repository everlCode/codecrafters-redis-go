package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type EchoCommand struct {
}

func (c EchoCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) == 0 {
		return Response(resp.Bulk("PONG"))
	}

	return Response(resp.Bulk(args[0]))
}
