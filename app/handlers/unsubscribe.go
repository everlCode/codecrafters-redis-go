package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type UnsubscribeCommand struct {
}

func (c UnsubscribeCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	chanelName := args[0]
	ps := server.Pubsub
	channel := ps.GetChannel(chanelName)
	if channel == nil {
		channel = pubsub.NewChannel(chanelName)
	}

	channel.Unbscribe(client.GetConnection())
	client.Unbscribe(chanelName)

	return Response(resp.Array([]any{"unsubscribe", chanelName, len(client.GetSubscribtions())}))
}
