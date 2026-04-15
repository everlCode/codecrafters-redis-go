package handlers

import (
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

	start := startValue.Integer
	end := endValue.Integer
	if (start > end) {
		return resp.Value{
			Type: resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	if (startValue.Type != resp.INTEGER || endValue.Type != resp.INTEGER) {
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

	v := value.Array[start : end + 1]

	return resp.Value{
		Type: resp.ARRAY,
		Array: v,
	}
}

func (c LRangeCommand) Name() string {
	return LRANGE
}
