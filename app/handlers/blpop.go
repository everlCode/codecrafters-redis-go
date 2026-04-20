package handlers

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type BlPopCommand struct {
}

func (c BlPopCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	key := args[0]
	if key.Type != resp.BULK {
		return resp.Value{Type: resp.ERROR, String: "Key should be string!"}
	}

	var timeout time.Duration = 0
	var endDate time.Time
	if len(args) > 1 {
		v, err := strconv.ParseFloat(args[1].Bulk, 64)
		if err != nil {
			return resp.Value{Type: resp.ERROR, String: err.Error()}
		}
		if v != 0 {
			timeout = time.Duration(v * float64(time.Second))
			endDate = time.Now().Add(timeout)
		}
	}

	entry, ok := db.Get(key.Bulk)
	if ok {
		var response resp.Value
		if entry.IsArray() {
			value, ok := entry.AsArray()
			if ok {
				if len(value) > 0 {
					firstValue := value[0]
					entry.Set(value[1:])
					db.Set(key.Bulk, entry)
					response = resp.Value{
						Type:  resp.ARRAY,
						Array: []resp.Value{resp.Value{Type: resp.BULK, Bulk: key.Bulk}, resp.Value{Type: resp.BULK, Bulk: firstValue}},
					}
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

	ch := make(chan database.Entry)
	db.PushWaiter(key.Bulk, &database.Waiter{Chanel: ch, Timeout: endDate})

	var response resp.Value

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

	return resp.Value{
		Type:  resp.ARRAY,
		Array: []resp.Value{resp.Value{Type: resp.BULK, Bulk: key.Bulk}, resp.Value{Type: resp.BULK, Bulk: response.Array[0].Bulk}},
	}
}
