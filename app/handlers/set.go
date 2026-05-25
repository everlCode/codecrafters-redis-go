package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type SetCommand struct {
}

func (c SetCommand) Execute(args []string, server *server.Server, client *clients.Client) resp.Value {
	db := server.GetDB()
	if len(args) < 2 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := args[0]
	value := args[1]
	var entry database.Entry
	entry.Set(value)

	if len(args) > 3 {
		commnadOption := args[2]
		if strings.ToLower(commnadOption) == "px" {
			px, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return resp.Value{Type: resp.ERROR, String: "Invalid argument"}
			}
			entry.Expires = time.Now().UnixMilli() + px
		}
	}

	db.Set(key, entry)

	return resp.SimpleString("OK")
}
