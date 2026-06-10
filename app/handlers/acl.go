package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type AclCommand struct {
}

func (c AclCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	_type := args[0]

	switch strings.ToLower(_type) {
	case "whoami":
		return Response(resp.Bulk("default"))
	case "getuser":
		return Response(resp.Array(
			[]any{
				"flags", []any{},
			},
		))
	}


    return Response(resp.Bulk("default"))
}