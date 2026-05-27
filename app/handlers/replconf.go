package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ReplconfCommand struct {
}

func (c ReplconfCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	var firstArgeument, secondArg string
	if len(args) > 0 {
		firstArgeument = args[0]
	}

	var ackCount int
	if len(args) > 1 {
		secondArg = args[1]
		secondArg, err := strconv.Atoi(secondArg)
		ackCount = secondArg
		if err != nil {
			return Response(resp.Error(err.Error()))
		}
	}
	
	if (firstArgeument == "GETACK") {
		offset :=  client.GetOffset()
		offetDigit := strconv.Itoa(offset)
		response := Response(resp.Array([]any{"REPLCONF", "ACK", offetDigit}))
		response.NeedAnswer = true

		return response
	}

	if firstArgeument == "ACK" && client.IsReplica() {
		replica := server.FindReplicaByClient(client)
		if replica != nil {
			replica.SetOffset(ackCount)	
		}
	}

	return Response(resp.SimpleString("OK"))
}
