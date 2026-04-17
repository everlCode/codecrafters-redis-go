package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type PingCommand struct {
}

func (c PingCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	return resp.Value{Type: resp.STRING, String: "PONG"}
}

func (c PingCommand) Name() string {
	return PING
}
