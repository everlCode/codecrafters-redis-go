package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type XrangeCommand struct {
}

func (c XrangeCommand) Execute(args []string, server *server.Server, client *clients.Client) resp.Value {
	db := server.GetDB()
	if len(args) < 3 {
		return resp.Error("ERR to few args!")
	}

	key := args[0]
	entry, ok := db.Get(key)
	if !ok {
		return resp.EmptyArray()
	}
	start := args[1]
	end := args[2]

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
