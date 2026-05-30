package dispatcher

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

var pubsubHandlers map[string]int = map[string]int{
	"SUBSCRIBE": 0,
	"UNSUBSCRIBE": 0,
	"PSUBSCRIBE": 0,
	"PUNSUBSCRIBE": 0,
	"PING": 0,
	"QUIT": 0,
}

func Dispatch(args []string, server *server.Server, client *clients.Client) handlers.CommandResponse {
	register := handlers.NewRegister()
	handlerName := strings.ToUpper(args[0])

	if client.IsSubscriber() && !isPubsubHandler(handlerName) {
		return handlers.Response(resp.Error("ERR Can't execute '" + handlerName + "': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context"))
	}

	switch handlerName {
	case "MULTI":
		client.StartTransactions()
		return handlers.Response(resp.SimpleString("OK"))
	case "EXEC":
		if !client.IsTransaction() {
			return handlers.Response(resp.Error("ERR EXEC without MULTI"))
		}
		client.EndTransactions()
		var execResp []any
		for _, c := range client.GetCommandQueue() {
			result := Dispatch(c.Args, server, client)
			execResp = append(execResp, result.GetResponse())
		}

		return handlers.Response(resp.Array(execResp))
	case "DISCARD":
		if !client.IsTransaction() {
			return handlers.Response(resp.Error("ERR DISCARD without MULTI"))
		}
		client.ClearCommandQueue()
		client.EndTransactions()

		return handlers.Response(resp.SimpleString("OK"))
	default:
		handler, err := register.Get(handlerName)

		if err != nil {
			return handlers.Response(resp.Error(fmt.Sprintf("ERR unknown command '%s'", handlerName)))
		}

		var result handlers.CommandResponse
		if client.IsTransaction() {
			commandQueue := clients.NewCommandQueue(handlerName, args)
			client.PushCommandQueue(commandQueue)

			result = handlers.Response(resp.SimpleString("QUEUED"))
		} else {
			result = handler.Execute(args[1:], server, client)
		}

		return result
	}
}

func isPubsubHandler(handlerName string) bool {
	_, ok := pubsubHandlers[handlerName]

	return ok
}