package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		fmt.Println(err)
		os.Exit(1)
	}

	db := database.New()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handle(conn, db)
	}
}

func handle(conn net.Conn, db *database.DB) {
	defer conn.Close()

	parser := resp.New(conn)
	writer := NewWriter(conn)

	for {
		request, err := parser.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			writer.Write(resp.Error(
				fmt.Sprintf("ERR: %s", err.Error()),
			))
			continue
		}
		if request.Type != resp.ARRAY {
			writer.Write(resp.Error(
				fmt.Sprintf("ERR: %s", "Неверный запрос! Ожидался массив!"),
			))
			continue
		}
		if len(request.Array) == 0 {
			writer.Write(resp.Error(
				fmt.Sprintf("ERR: %s", "Неверный запрос! Не переданы необходимые аргументы"),
			))

			continue
		}

		register := handlers.NewRegister()
		handlerName := strings.ToUpper(request.Array[0].Bulk)

		handler, err := register.Get(handlerName)

		if err != nil {
			writer.Write(resp.Error(
				fmt.Sprintf("ERR unknown command '%s'", handlerName),
			))
			continue
		}

		args := request.Array[1:]
		result := handler.Execute(args, db)

		writer.Write(result)
	}
}
