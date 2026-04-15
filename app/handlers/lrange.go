package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LRangeCommand struct {
}

func (c LRangeCommand) Execute(args []resp.Value, db *db.DB) resp.Value {
	if len(args) < 3 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}
	key := args[0]
	startValue := args[1]
	endValue := args[2]

	if (startValue.Type != resp.BULK || endValue.Type != resp.BULK) {
		return resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	start, err := strconv.Atoi(startValue.Bulk)
	if err != nil {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}
	end, err := strconv.Atoi(endValue.Bulk)
	if err != nil {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}
	
	if (start > end) {
		return resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	value, ok := db.Get(key.Bulk)
	if !ok || value.Type != resp.ARRAY{
		return resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	lenght := len(value.Array)

	if start >= lenght {
		return resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	if end > lenght {
		end = lenght
	}

	v := value.Array[start : end]

	return resp.Value{
		Type: resp.ARRAY,
		Array: v,
	}
}

func (c LRangeCommand) Name() string {
	return LRANGE
}
