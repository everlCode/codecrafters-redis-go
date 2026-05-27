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
	var firstArgeument string
	if len(args) > 0 {
		firstArgeument = args[0]
	}

	if (firstArgeument == "GETACK") {
		return c.getACK(client)
	}

	if firstArgeument == "ACK" && client.IsReplica() {
		var ackCount int
		if len(args) > 1 {
			secondArg := args[1]
			v, err := strconv.Atoi(secondArg)
			if err != nil {
				return Response(resp.Error(err.Error()))
			}
			ackCount = v
		}
		c.ack(ackCount, client, server)
	}

	return Response(resp.SimpleString("OK"))
}

func (c ReplconfCommand) getACK(client *clients.Client) CommandResponse {
	offset :=  client.GetOffset()
	offetDigit := strconv.Itoa(offset)
	response := Response(resp.Array([]any{"REPLCONF", "ACK", offetDigit}))
	response.NeedAnswer = true

	return response
}

func (c ReplconfCommand) ack(count int, client *clients.Client, server *server.Server) {
	
	replica := server.FindReplicaByClient(client)
	if replica != nil {
		replica.SetOffset(count)	
	}
}
