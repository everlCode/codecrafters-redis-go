package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type BlPopCommand struct {
}

func (c BlPopCommand) Execute(args []string, db *database.DB) resp.Value {
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
				response = resp.Value{
					Type:  resp.ARRAY,
					Array: []resp.Value{resp.Value{Type: resp.BULK, Bulk: key}, resp.Value{Type: resp.BULK, Bulk: firstValue}},
				}
			}
		} else {
			response = resp.Value{
				Type:  resp.ARRAY,
				Array: nil,
			}
		}

		return response
	}

	ch := db.PushWaiter(key, endDate)

	if timeout == 0 {
		entry = <-ch
	} else {
		select {
		case v := <-ch:
			entry = v
		case <-time.After(timeout):
			return resp.Value{
				Type:  resp.ARRAY,
				Array: nil,
			}
		}
	}
	data := entry.AsArray()
	if len(data) > 0 {
		data = []string{key, data[0]}
	}

	return resp.ArrayString(data)
}
