package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type LRangeCommand struct {
}

func (c LRangeCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
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
	if err != nil {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}
	end, err := strconv.Atoi(endValue.Bulk)
	if err != nil {
		return resp.Value{Type: resp.ERROR, Bulk: err.Error()}
	}

	entry, ok := db.Get(key.Bulk)
	data, _ := entry.AsArray()
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

	return resp.Array(data)
}