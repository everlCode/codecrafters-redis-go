package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type SubscribeCommand struct {
}

func (c SubscribeCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	chanelName := args[0]
	ps := server.Pubsub
	var channel *pubsub.Channel
	channel = ps.GetChannel(chanelName)
	if channel == nil {
		channel = pubsub.NewChannel(chanelName)
	}

	ps.Subscribe(channel, client.GetConnection())
	client.Subscribe(channel)
	
	

	return Response(resp.Array([]any{"subscribe", chanelName, len(client.GetSubscribtions())}))
}
