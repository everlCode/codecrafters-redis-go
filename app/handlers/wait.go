package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type WaitCommand struct {
}

func (c WaitCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	var count int = 0
	if len(server.Replicas) > 0 {
		count = len(server.Replicas)
	}

	return Response(resp.Integer(count))
}
