package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type BlPopCommand struct {
}

func (c BlPopCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	db := server.GetDB()
	key := args[0]

	var argParser = New()
	var timeout time.Duration = 0
	var endDate time.Time
	if len(args) > 1 {
		timeout = argParser.ParseTimeout(args[1], true)
	}
	if timeout > 0 {
		endDate = time.Now().Add(timeout)
	}
	
	entry, ok := db.Get(key)
	if ok {
		var response resp.Value
		if entry.IsArray() {
			value := entry.AsArray()
			if len(value) > 0 {
				firstValue := value[0]
				entry.Set(value[1:])
				db.Set(key, entry)

				response = resp.ArrayString([]string{key, firstValue})
			}
		} else {
			response = resp.NilArray()
		}

		return Response(response)
	}

	ch := db.PushWaiter(key, endDate)

	if timeout == 0 {
		entry = <-ch
	} else {
		select {
		case v := <-ch:
			entry = v
		case <-time.After(timeout):
			return Response(resp.NilArray())
		}
	}
	data := entry.AsArray()
	if len(data) > 0 {
		data = []string{key, data[0]}
	}

	return Response(resp.ArrayString(data))
}
