package main

import "fmt"

var handlers = map[string]func([]Value, *DB) Value{
	"echo": echo,
	"set": set,
	"get": get,
}

func echo(args []Value, db *DB) Value {
	fmt.Println("ECHO")
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

	key := args[0]
	value := args[1]

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
	fmt.Println(db.sets)
	
	if !ok {
		return Value{Type: ERROR, String: ""}
	}

	return value
}
