package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	db := NewDB()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handle(conn, db)
	}
}

func handle(conn net.Conn, db *DB) {
	defer conn.Close()

	for {
		parser := NewResp(conn)
		command, err := parser.Read()
		if err != nil {
			fmt.Errorf("ERR: %s", err.Error())
			continue
		}
		if command.Type != ARRAY {
			fmt.Errorf("ERR: %s", "Неверный запрос! Ожидался массив!")
			continue
		}
		if len(command.Array) == 0 {
			fmt.Errorf("ERR: %s", "Неверный запрос! Не переданы необходимые аргументы")
			continue
		}

		handlerName := command.Array[0].Bulk
		handler, ok := handlers[strings.ToLower(handlerName)]
		if !ok {
			fmt.Errorf("ERR: %s", "Неверный запрос! Команды %s не существует!", handlerName)
			continue
		}
		r := command.Array[1:]
		result := handler(r, db)

		writer := NewWriter(conn)
		writer.Write(result)
	}
}
