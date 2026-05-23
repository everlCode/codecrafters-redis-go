package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type MultiCommand struct {
}

func (c MultiCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	return resp.SimpleString("OK")
}
