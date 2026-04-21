package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type XaddCommand struct {
}

func (c XaddCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	argss := resp.ParseSlice(args)
	if len(argss) < 4 {
		return resp.Error("ERR to few args!")
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
	var stream database.Stream
	if ok {
		stream, _ = entry.AsStream()

		_, err := c.validateId(id, stream)
		if err != nil {
			return resp.Error(err.Error())
		}

		stream.Add(id, streamData)
		entry.Set(stream)
	} else {
		stream = database.CreateStream(id, streamData)
		entry = database.Entry{}
		entry.Set(stream)
	}

	db.Set(key, entry)

	return resp.Bulk(id)
}

func (c XaddCommand) validateId(id string, stream database.Stream) (bool, error) {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return false, nil
	}

	miliseconds, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, nil
	}
	indetificator, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, nil
	}
	if indetificator == 0 {
		return false, errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}

	lastIdUnixTime, lastId := stream.GetLastIdParts()

	if miliseconds < lastIdUnixTime {
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if miliseconds == lastIdUnixTime && indetificator < lastId {
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if miliseconds == lastIdUnixTime && indetificator == lastId {
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return true, nil
}
