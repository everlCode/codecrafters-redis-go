package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type SetCommand struct {
}

func (c SetCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	fmt.Println("SET")
	if len(args) < 2 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}
	fmt.Println(args)

	key := args[0]
	value := args[1]

	if len(args) > 3 {
		commnadOption := args[2]
		if strings.ToLower(commnadOption.Bulk) == "px" {
			px, err := strconv.ParseInt(args[3].Bulk, 10, 64)
			if err != nil {
				return resp.Value{Type: resp.ERROR, String: "Invalid argument"}
			}
			value.Expires = time.Now().UnixMilli() + px
		}
	}

	db.Sets[key.Bulk] = value

	return resp.Value{Type: resp.STRING, String: "OK"}
}

func (c SetCommand) Name() string {
	return SET
}
