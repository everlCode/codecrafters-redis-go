package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ExecCommand struct {
}

func (c ExecCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	if !db.IsTransaction() {
		return resp.Error("ERR EXEC without MULTI")
	} 

	return resp.SimpleString("OK")
}
