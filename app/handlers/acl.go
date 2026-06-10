package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type AclCommand struct {
}

func (c AclCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
    return Response(resp.Bulk("default"))
}