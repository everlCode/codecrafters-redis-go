package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

const (
	PING        = "PING"
	SET         = "SET"
	GET         = "GET"
	ECHO        = "ECHO"
	COMMAND     = "COMMAND"
	LPUSH       = "LPUSH"
	RPUSH       = "RPUSH"
	LRANGE      = "LRANGE"
	LLEN        = "LLEN"
	LPOP        = "LPOP"
	BLPOP       = "BLPOP"
	TYPE        = "TYPE"
	XADD        = "XADD"
	XRANGE      = "XRANGE"
	XREAD       = "XREAD"
	INCR        = "INCR"
	INFO        = "INFO"
	REPLCONF    = "REPLCONF"
	PSYNC       = "PSYNC"
	WAIT        = "WAIT"
	CONFIG      = "CONFIG"
	KEYS        = "KEYS"
	SUBSCRIBE   = "SUBSCRIBE"
	PUBLISH     = "PUBLISH"
	UNSUBSCRIBE = "UNSUBSCRIBE"
	ZADD        = "ZADD"
	ZRANK       = "ZRANK"
	ZRANGE      = "ZRANGE"
	ZCARD      = "ZCARD"
	ZSCORE       = "ZSCORE"
	ZREM       = "ZREM"
	GEOADD       = "GEOADD"
	GEOPOS       = "GEOPOS"
	GEODIST       = "GEODIST"
	GEOSEARCH       = "GEOSEARCH"
	ACL       = "ACL"
	AUTH       = "AUTH"
)

type Command interface {
	Execute(
		args []string,
		server *server.Server,
		client *clients.Client,
	) CommandResponse
}

type CommandResponse struct {
	value         resp.Value
	IsPropagation bool
	NeedAnswer    bool
	PropCommand   PropagationCommand
}

type PropagationCommand struct {
	Name string
	Args []string
}

func (cr *CommandResponse) GetResponse() resp.Value {
	return cr.value
}

func (cr *CommandResponse) SetPropagationCommand(name string, args []string) {
	cr.IsPropagation = true
	cr.PropCommand = PropagationCommand{
		Name: name,
		Args: args,
	}
}

func Response(data resp.Value) CommandResponse {
	return CommandResponse{
		value: data,
	}
}
