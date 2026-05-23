package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type MultiCommand struct {
}

func (c MultiCommand) Execute(args []string, db *database.DB) resp.Value {
	db.SetMulti(true)

	return resp.SimpleString("OK")
}
