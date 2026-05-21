package handlers

import (
	"strings"
	"time"

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

	parsedArgs := c.parseArgs(argss)
	isBlock := parsedArgs["isBlock"].(bool)
	streamKeyMap := parsedArgs["pairs"].(map[string]string)
	timeout := parsedArgs["timeout"].(time.Duration)

	var endDate time.Time
	endDate = time.Now().Add(timeout)

	var response []any
	for key, start := range streamKeyMap {
		entry, _ := db.Get(key)
		stream, _ := entry.AsStream()
		streamEntries := stream.GetEntries(start)
		if len(streamEntries) > 0 {
			var streamEntryData []any = []any{}
			for _, streamEntry := range streamEntries {
				streamEntryData = append(streamEntryData, PrepareStreamEntryData(streamEntry))

			}
			streamData := []any{key, streamEntryData}
			response = append(response, streamData)
		} else if isBlock {
			ch := db.PushWaiter(key, endDate)

			if timeout == 0 {
				entry = <-ch
			} else {
				select {
				case v := <-ch:
					entry = v
				case <-time.After(timeout):
					return resp.Value{
						Type:  resp.ARRAY,
						Array: nil,
					}
				}
			}
		} else {
			return resp.EmptyArray()
		}
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

	if isBlock {
		var argParser = New()
		var timeout time.Duration = 0
		if len(args) > 1 {
			timeout = argParser.ParseTimeout(arguments[1])
		}
		response["timeout"] = timeout
		pairs = arguments[2:]
	} else {
		pairs = arguments
	}
	response["isBlock"] = isBlock

	pairsLen := len(pairs)
	pairsMiddle := pairsLen / 2
	var streamKeyMap map[string]string = make(map[string]string)
	for i, v := range pairs {
		if i+1 > pairsMiddle {
			break
		}
		streamKeyMap[v] = pairs[i+pairsMiddle]
	}
	response["pairs"] = streamKeyMap

	return response
}
