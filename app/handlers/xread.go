package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type XreadCommand struct {
}

func (c XreadCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	argss := resp.ParseSlice(args)
	if len(argss) < 3 {
		return resp.Error("ERR to few args!")
	}

	pairs := argss[1:]
	pairsLen := len(pairs)
	pairsMiddle := pairsLen / 2
	var streamKeyMap map[string]string = make(map[string]string)
	for i, v := range pairs {
		if i + 1 > pairsMiddle {
			break
		}
		streamKeyMap[v] = pairs[i + pairsMiddle]
	}

	var response []any
	for key, start := range streamKeyMap {
		entry, ok := db.Get(key)
		if !ok {
			return resp.EmptyArray()
		}

		stream, _ := entry.AsStream()
		streamEntries := stream.GetEntries(start)

		var streamEntryData []any = []any{}
		for _, streamEntry := range streamEntries {
			streamEntryData = append(streamEntryData, PrepareStreamEntryData(streamEntry))

		}
		streamData := []any{key, streamEntryData}
		response = append(response, streamData)
	}
	
	return resp.Array(response)
}
