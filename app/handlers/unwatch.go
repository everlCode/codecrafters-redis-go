package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type UnwatchCommand struct {
}

func (c UnwatchCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	client.Unwatch()
	

	return Response(resp.SimpleString("OK"))
}
