package handlers

import (
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type XreadCommand struct {
}

func (c XreadCommand) Execute(args []resp.Value, db *database.DB) resp.Value {
	argss := resp.ParseSlice(args)
	if len(argss) < 3 {
		return resp.Error("ERR to few args!")
	}

	parsedArgs := c.parseArgs(argss)
	isBlock := parsedArgs["isBlock"].(bool)
	streamKeys := parsedArgs["keys"].([]string)
	starts := parsedArgs["starts"].([]string)
	timeout := parsedArgs["timeout"].(time.Duration)

	var endDate time.Time
	if (timeout > 0) {
		endDate = time.Now().Add(timeout)
	}

	var response []any
	for i, key := range streamKeys {
		var entry database.Entry
		var stream *database.Stream
		var streamEntriesCountBefore int

		start := starts[i]

		entry, _ = db.Get(key)
		if isBlock {
			if start == "$" {
				stream, _ = entry.AsStream()
				streamEntriesCountBefore = len(stream.GetEntries())
			}
			ch := db.PushWaiter(key, endDate)
			var timeoutCh <-chan time.Time
			if timeout > 0 {
				timeoutCh = time.After(timeout)
			}

			select {
			case v := <-ch:
				entry = v

			case <-timeoutCh:
				return resp.Value{
					Type:  resp.ARRAY,
					Array: nil,
				}
			}
		}
		stream, _ = entry.AsStream()
		startvalue := helpers.IncrementStreamId(start)
		streamEntries := stream.GetEntries(startvalue)
		var streamEntryData []any = []any{}
		for k, streamEntry := range streamEntries {
			if start == "$" && k < streamEntriesCountBefore {
				continue
			}
			streamEntryData = append(streamEntryData, PrepareStreamEntryData(streamEntry))

		}
		streamData := []any{key, streamEntryData}
		response = append(response, streamData)
	}

	return resp.Array(response)
}

func (c XreadCommand) parseArgs(args []string) map[string]any {
	arguments := make([]string, 0)
	for _, v := range args {
		if strings.ToLower(v) != "streams" {
			arguments = append(arguments, v)
		}
	}
	isBlock := strings.ToLower(arguments[0]) == "block"
	var pairs []string
	var response map[string]any = make(map[string]any)
	var timeout time.Duration = 0
	response["timeout"] = timeout

	if isBlock {
		var argParser = New()
		var timeout time.Duration = 0
		if len(args) > 1 {
			timeout = argParser.ParseTimeout(arguments[1], false)
		}
		response["timeout"] = timeout
		pairs = arguments[2:]
	} else {
		pairs = arguments
	}
	response["isBlock"] = isBlock

	pairsLen := len(pairs)
	pairsMiddle := pairsLen / 2
	
	var starts []string
	var keys []string
	for i, v := range pairs {
		if i+1 > pairsMiddle {
			starts = append(starts, v)
		} else {
			keys = append(keys, v)
		}
		
	}
	response["keys"] = keys
	response["starts"] = starts

	return response
}
