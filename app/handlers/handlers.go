package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	PING    = "PING"
	SET     = "SET"
	GET     = "GET"
	ECHO    = "ECHO"
	COMMAND = "COMMAND"
	LPUSH   = "LPUSH"
	RPUSH   = "RPUSH"
	LRANGE  = "LRANGE"
	LLEN    = "LLEN"
	LPOP    = "LPOP"
	BLPOP   = "BLPOP"
	TYPE    = "TYPE"
	XADD    = "XADD"
	XRANGE  = "XRANGE"
	XREAD   = "XREAD"
	INCR   = "INCR"
	INFO   = "INFO"
)

type Command interface {
	Execute(args []string, db *database.DB, client *clients.Client) resp.Value
}
