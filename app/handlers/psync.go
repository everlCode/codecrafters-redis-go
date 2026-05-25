package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type PsyncCommand struct {
}

func (c PsyncCommand) Execute(args []string, server *server.Server) resp.Value {


	return resp.SimpleString("FULLRESYNC " + server.MasterReplyId + " 0")
}
