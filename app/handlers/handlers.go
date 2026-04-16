package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	PING = "PING"
	SET  = "SET"
	GET  = "GET"
	ECHO = "ECHO"
	COMMAND = "COMMAND"
	LPUSH = "LPUSH"
	RPUSH = "RPUSH"
	LRANGE = "LRANGE"
	LLEN = "LLEN"
)

type Command interface {
	Execute([]resp.Value, *db.DB) resp.Value
	Name() string
}