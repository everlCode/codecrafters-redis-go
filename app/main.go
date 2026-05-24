package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/dispatcher"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	db := database.New()
	server := server.New(db)
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:" + server.Port)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handle(conn, server)
	}
}

func handle(conn net.Conn, server *server.Server) {
	defer conn.Close()

	parser := resp.New(conn)
	writer := NewWriter(conn)
	client := clients.New()

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
		args := resp.ParseSlice(request.Array)
		result := dispatcher.Dispatch(args, server, client)
		writer.Write(result)
	}
}

