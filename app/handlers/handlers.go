package handlers

import (
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
	MULTI   = "MULTI"
	EXEC   = "EXEC"
)

type Command interface {
	Execute([]string, *database.DB) resp.Value
}
