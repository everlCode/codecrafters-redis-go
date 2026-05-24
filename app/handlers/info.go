package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type InfoCommand struct {
}

func (c InfoCommand) Execute(args []string, server *server.Server) resp.Value {
	key := args[0]
	config := server.GetConfig()
	
	switch strings.ToLower(key) {
	case "replication":
		return resp.Bulk("# Replication \n role:" + config.Role)
	}

	return resp.Bulk("")
}
