package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type LRangeCommand struct {
}

func (c LRangeCommand) Execute(args []string, server *server.Server, client *clients.Client) resp.Value {
	db := server.GetDB()
	if len(args) < 3 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}
	key := args[0]
	startValue := args[1]
	endValue := args[2]

	start, err := strconv.Atoi(startValue)
	if err != nil {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}
	end, err := strconv.Atoi(endValue)
	if err != nil {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}

	entry, ok := db.Get(key)
	data := entry.AsArray()
	lenght := len(data)

	if start < 0 {
		start = lenght + start
		if start < 0 {
			start = 0
		}
	}

	if end < -1 {
		end = lenght + end + 1
	} else if end < lenght && end > 0 {
		end += 1
	}

	if start > end && end >= 0 {
		return resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	if !ok || !entry.IsArray() {
		return resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	if start >= lenght {
		return resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	if end == -1 || end >= lenght || end == 0 {
		data = data[start:]
	} else {
		data = data[start:end]
	}

	return resp.ArrayString(data)
}