package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type XrangeCommand struct {
}

func (c XrangeCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	argss := resp.ParseSlice(args)
	if len(argss) < 3 {
		return resp.Error("ERR to few args!")
	}

	key := argss[0]
	entry, ok := db.Get(key)
	if !ok {
		return resp.EmptyArray()
	}
	start := argss[1]
	end := argss[2]

	stream, _ := entry.AsStream()
	streamEntries := stream.GetEntries(start, end)

	var response []any
	for _, streamEntry := range streamEntries {
		response = append(response, PrepareStreamEntryData(streamEntry))
	}

	return resp.Array(response)
}

func PrepareStreamEntryData(streamEntry database.StreamEntry) []any {
	entryId := streamEntry.GetId()
	dataMap := streamEntry.GeData()
	var preparedData []any
	for key, value := range dataMap {
		preparedData = append(preparedData, key, value)
	}

	return []any{entryId, preparedData}
}
