package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type PublishCommand struct {
}

func (c PublishCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 2 {
		return Response(resp.Error("Too few args!"))
	}
	channelName := args[0]
	content := args[1]
	channel := server.Pubsub.GetChannel(channelName)

	for _, v := range channel.Connections {
		server.SendRequest(v, "message", channelName, content)
	}

	return Response(resp.Integer(len(channel.Connections)))
}
