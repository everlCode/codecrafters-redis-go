package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
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
	var stream *database.Stream
	var lastId string
	if ok {
		stream, _ = entry.AsStream()
		lastId = stream.GetLastId()

		entry.Set(stream)
	} else {
		stream = database.CreateStream()
		entry = database.Entry{}
		entry.Set(stream)
	}
	generatedId := c.generateId(id, lastId)
	_, err := c.validateId(generatedId, lastId)
	if err != nil {
		return resp.Error(err.Error())
	}
	stream.Add(generatedId, streamData)

	db.Set(key, entry)

	return resp.Bulk(generatedId)
}

func (c XaddCommand) validateId(id string, lastId string) (bool, error) {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return false, nil
	}

	miliseconds, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, nil
	}

	lastIdUnixTime, lastIndetificator := helpers.GetStreamIdParts(lastId)
	indetificator, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, nil
	}
	if indetificator == 0 && miliseconds == 0 {
		return false, errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	if miliseconds < lastIdUnixTime {
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if miliseconds == lastIdUnixTime && indetificator < lastIndetificator {
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if miliseconds == lastIdUnixTime && indetificator == lastIndetificator {
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return true, nil
}

func (c XaddCommand) generateId(id string, lastId string) string {
	var parts []string
	if id == "*" {
		parts = []string{"*", "*"}
	} else {
		parts = strings.Split(id, "-")
		if len(parts) < 2 {
			return "0-1"
		}
	}
	if parts[0] == "*" {
		curTime := time.Now().UnixMilli()
		parts[0] = strconv.Itoa(int(curTime))
	}
	lastTime, lastIndetificator := helpers.GetStreamIdParts(lastId)
	currentTime, _ := strconv.Atoi(parts[0])
	if parts[1] == "*" {
		if lastTime == currentTime {
			parts[1] = strconv.Itoa(lastIndetificator + 1)
		} else {
			parts[1] = "0"
		}
	}
	str := strings.Join(parts, "-")

	return str
}
