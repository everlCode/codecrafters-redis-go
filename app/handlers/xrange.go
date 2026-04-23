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
	start := argss[1]
	var startMilliseconds, startSeqId int
	if start == "-" {
		startMilliseconds, startSeqId = 0, 0
	} else {
		startMilliseconds, startSeqId = GetStreamIdParts(start)
	}

	end := argss[2]
	var endMilliseconds, endSeqId int
	isEndExist := end != "+"
	if isEndExist {
		endMilliseconds, endSeqId = GetStreamIdParts(end)
	}

	entry, ok := db.Get(key)
	if !ok {
		return resp.EmptyArray()
	}

	stream, _ := entry.AsStream()
	streamEnties := stream.GetEntries()

	var response []any
	for _, streamEntry := range streamEnties {
		entryId := streamEntry.GetId()
		miliseconds, seqId := GetStreamIdParts(entryId)

		if miliseconds >= startMilliseconds && seqId >= startSeqId &&
			(!isEndExist || (miliseconds <= endMilliseconds && seqId <= endSeqId)) {

			dataMap := streamEntry.GeData()
			var preparedData []any
			for key, value := range dataMap {
				preparedData = append(preparedData, key, value)
			}
			iddd := streamEntry.GetId()

			response = append(response, []any{iddd, preparedData})
		}
	}

	return resp.Array(response)
}
