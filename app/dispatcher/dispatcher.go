package dispatcher

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func Dispatch(args []string, server *server.Server, client *clients.Client) resp.Value {
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
			execResp = append(execResp, Dispatch(c.Args, server, client))
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
			result = handler.Execute(args[1:], server, client)
		}

		return result
	}
}