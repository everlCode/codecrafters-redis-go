package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/acl"
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
		user := server.Acl.GetUser("default")

		return generateResponse(*user)
	case "setuser":
		return c.setUser(args[0:], server)
	}

	return Response(resp.Bulk("default"))
}

func generateResponse(user acl.User) CommandResponse {
	flags := []any{}
	passwords := []any{}
	if user.NoPass {
		flags = append(flags, "nopass")
	} else {
		passwords = append(passwords, user.Password)
	}

	return Response(
		resp.Array(
			[]any{
				"flags", flags, "passwords", passwords,
			}),
	)

}

func (c AclCommand) setUser(args []string, server *server.Server) CommandResponse {
	userName := args[1]
	user := server.Acl.GetUser(userName)
	if user == nil {
		return Response(resp.Error("User not found!"))
	}

	param := args[2]
	switch string(param[0]) {
	case ">":
		user.SetPassword(string(param[1:]))
	}

	return Response(resp.SimpleString("OK"))
}
