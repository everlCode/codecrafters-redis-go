package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type WatchCommand struct {
}

func (c WatchCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	return Response(resp.SimpleString("OK"))
}
