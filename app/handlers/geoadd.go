package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GeoaddCommand struct {
}

func (c GeoaddCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	
	return Response(resp.Integer(1))
}
