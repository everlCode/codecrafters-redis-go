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

		timeout = time.Duration(v * float64(time.Second))
		endDate = time.Now().Add(timeout)
	}

	value, ok := db.Get(key.Bulk)
	if ok {
		v := value.Array[0]
		value.Array = value.Array[1:]

		db.Set(key.Bulk, value)

		value = resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{resp.Value{Type: resp.BULK, Bulk: key.Bulk}, resp.Value{Type: resp.BULK, Bulk: v.Bulk}},
		}

		return value
	}

	ch := make(chan resp.Value)
	db.PushWaiter(key.Bulk, &database.Waiter{Chanel: ch, Timeout: endDate})

	var response resp.Value
	select {
	case v := <-ch:
		response = v
	case <-time.After(timeout):
		return resp.Value{
			Type:  resp.ARRAY,
			Array: nil,
		}
	}

	return resp.Value{
		Type:  resp.ARRAY,
		Array: []resp.Value{resp.Value{Type: resp.BULK, Bulk: key.Bulk}, resp.Value{Type: resp.BULK, Bulk: response.Array[0].Bulk}},
	}
}

func (c BlPopCommand) Name() string {
	return BLPOP
}
