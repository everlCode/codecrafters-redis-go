package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type SetCommand struct {
}

func (c SetCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	if len(args) < 2 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := args[0]
	value := args[1]
	var entry database.Entry
	entry.Set(value.Marshal())

	if len(args) > 3 {
		commnadOption := args[2]
		if strings.ToLower(commnadOption.Bulk) == "px" {
			px, err := strconv.ParseInt(args[3].Bulk, 10, 64)
			if err != nil {
				return resp.Value{Type: resp.ERROR, String: "Invalid argument"}
			}
			entry.Expires = time.Now().UnixMilli() + px
		}
	}

	db.Set(key.Bulk, entry)

	return resp.SimpleString("OK")
}
