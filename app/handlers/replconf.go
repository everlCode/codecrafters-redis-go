package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ReplconfCommand struct {
}

func (c ReplconfCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	var firstArgeument string
	if len(args) > 0 {
		firstArgeument = args[0]
	}
	if (firstArgeument == "GETACK") {
		response := Response(resp.Array([]any{"REPLCONF", "ACK", "0"}))
		response.NeedAnswer = true

		return response
	}

	return Response(resp.SimpleString("OK"))
}
