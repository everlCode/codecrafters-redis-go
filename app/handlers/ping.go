package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type PingCommand struct {
}

func (c PingCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if client.IsSubscriber() {
		return Response(resp.ArrayString("pong", ""))
	}

	return Response(resp.SimpleString("PONG"))
}
