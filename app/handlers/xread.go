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

	key := argss[1]
	entry, ok := db.Get(key)
	if !ok {
		return resp.EmptyArray()
	}
	start := argss[2]

	stream, _ := entry.AsStream()
	streamEntries := stream.GetEntries(start)

	var response []any
	var streamEntryData []any = []any{}
	for _, streamEntry := range streamEntries {
		streamEntryData = append(streamEntryData, PrepareStreamEntryData(streamEntry))

	}
	streamData := []any{key, streamEntryData}
	response = append(response, streamData)

	return resp.Array(response)
}
