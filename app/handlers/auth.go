package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type AuthCommand struct {
}

func (c AuthCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	username := args[0]
	user := server.Acl.GetUser(username)
	pass := args[1]
	if user == nil || !user.CheckPassword(pass) {
		return Response(resp.Error("WRONGPASS invalid username-password pair or user is disabled."))
	}
	client.User = user
	client.Authenticated = true

	
	return Response(resp.SimpleString("OK"))
}
