package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type InfoCommand struct {
}

func (c InfoCommand) Execute(args []string, server *server.Server, client *clients.Client) resp.Value {
	key := args[0]
	
	switch strings.ToLower(key) {
	case "replication":
		return resp.Bulk("# Replication \r\n role:" + server.Role + "\r\n" + 
		"master_replid:" + server.MasterReplyId + "\r\n" + 
		"master_repl_offset:" + server.MasterReplyOffset)
	}
	return resp.Bulk("")
}
