package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var handlers = map[string]func([]Value, *DB) Value{
	"ping": ping,
	"echo": echo,
	"set":  set,
	"get":  get,
}

func ping(args []Value, db *DB) Value {
	return Value{Type: STRING, String: "PONG"}
}

func echo(args []Value, db *DB) Value {
	if len(args) == 0 {
		return Value{Type: BULK, Bulk: "PONG"}
	}

	return Value{Type: BULK, Bulk: args[0].Bulk}
}

func set(args []Value, db *DB) Value {
	fmt.Println("SET")
	if len(args) < 2 {
		return Value{Type: ERROR, String: "ERR to few args"}
	}
	fmt.Println(args)

	key := args[0]
	value := args[1]

	if len(args) > 3 {
		commnadOption := args[2]
		if strings.ToLower(commnadOption.Bulk) == "px" {
			px, err := strconv.ParseInt(args[3].Bulk, 10, 64)
			if err != nil {
				return Value{Type: ERROR, String: "Invalid argument"}
			}
			value.Expires = time.Now().UnixMilli() + px
		}
	}

	db.sets[key.Bulk] = value

	return Value{Type: STRING, String: "OK"}
}

func get(args []Value, db *DB) Value {
	fmt.Println("GET")
	if len(args) < 1 {
		return Value{Type: ERROR, String: "ERR to few args"}
	}

	key := args[0]

	value, ok := db.sets[key.Bulk]
	fmt.Println("NOW ", time.Now().UnixMilli())
	fmt.Println("VALUE ", value.Expires)
	if !ok || (value.Expires != 0 && value.Expires < time.Now().UnixMilli()) {
		return Value{Type: BULK, Bulk: ""}
	}

	return value
}
