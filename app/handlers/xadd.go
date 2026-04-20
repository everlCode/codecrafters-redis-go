package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type XaddCommand struct {
}

func (c XaddCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	argss := resp.ParseSlice(args)
	if len(argss) < 4 {
		return resp.Value{Type: resp.ERROR, String: "ERR to few args"}
	}

	key := argss[0]
	id := argss[1]
	data := argss[2:]

	streamData := make(map[string]string)
	for i := 0; i < len(data); i++ {
		if i%2 == 0 {
			streamData[data[i]] = ""
		} else {
			streamData[data[i-1]] = data[i]
		}
	}

	entry, ok := db.Get(key)
	if !ok {
		entry = database.CreateStream(id, streamData)
	} else {
		stream, _ := entry.AsStream()
		stream.Add(id, streamData)
		entry.Set(stream)
	}

	return resp.Bulk(id)
}
