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

	if startValue.Type != resp.BULK || endValue.Type != resp.BULK {
		return resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	start, err := strconv.Atoi(startValue.Bulk)
	if (err != nil) {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}
	end, err := strconv.Atoi(endValue.Bulk)
	if (err != nil) {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}

	value, ok := db.Get(key.Bulk)
	lenght := len(value.Array)

	if (start < 0) {
		start = lenght + start
		if (start < 0) {
			start = 0
		}
	}

	if (end < -1) {
		end = lenght + end + 1
	} else if (end < lenght && end > 0){
		end += 1
	}

	if start > end && end > 0 {
		return resp.Value{
			Type:  resp.ARRAY,
			Array: []resp.Value{},
		}
	}

	if !ok || value.Type != resp.ARRAY {
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


	var v []resp.Value

	if (end == -1 || end >= lenght || end == 0) {
		v = value.Array[start:]
	} else {
		v = value.Array[start : end]
	}
	
	return resp.Value{
		Type:  resp.ARRAY,
		Array: v,
	}
}

func (c LRangeCommand) Name() string {
	return LRANGE
}
