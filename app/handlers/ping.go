package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type PingCommand struct {
}

func (c PingCommand) Execute(args []string, db *database.DB, client *clients.Client) resp.Value {
	return resp.Value{Type: resp.STRING, String: "PONG"}
}
