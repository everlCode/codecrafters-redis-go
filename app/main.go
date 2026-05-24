package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
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
		result := Dispatch(args, db, client)
		writer.Write(result)
	}
}

func Dispatch(args []string, db *database.DB, client *clients.Client) resp.Value {
	register := handlers.NewRegister()
	handlerName := strings.ToUpper(args[0])

	switch handlerName {
	case "MULTI":
		client.StartTransactions()
		return resp.SimpleString("OK")
	case "EXEC":
		if !client.IsTransaction() {
			return resp.Error("ERR EXEC without MULTI")
		}
		client.EndTransactions()
		var execResp []any
		for _, c := range client.GetCommandQueue() {
			execResp = append(execResp, Dispatch(c.Args, db, client))
		}

		return resp.Array(execResp)
	case "DISCARD":
		if !client.IsTransaction() {
			return resp.Error("ERR DISCARD without MULTI")
		}
		client.ClearCommandQueue()
		client.EndTransactions()

		return resp.SimpleString("OK")
	default:
		handler, err := register.Get(handlerName)

		if err != nil {
			return resp.Error(fmt.Sprintf("ERR unknown command '%s'", handlerName))
		}

		var result resp.Value
		if client.IsTransaction() {
			commandQueue := clients.NewCommandQueue(handlerName, args)
			client.PushCommandQueue(commandQueue)

			result = resp.SimpleString("QUEUED")
		} else {
			result = handler.Execute(args[1:], db, client)
		}

		return result
	}
}
