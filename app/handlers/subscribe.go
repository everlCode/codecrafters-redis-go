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
	subscription := pubsub.NewSubscription(client.GetConnection())
	ps.AddSubscription(subscription)

	client.Subscribe(subscription)
	subscriptions := client.GetSubscribtions()

	return Response(resp.Array([]any{"subscribe", chanelName, len(subscriptions)}))
}
