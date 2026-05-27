package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type WaitCommand struct {
}

func (c WaitCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	return Response(resp.Integer(0))
}
