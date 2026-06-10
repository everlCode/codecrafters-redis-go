package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GeosearchCommand struct {
}

func (c GeosearchCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
    if len(args) < 6 {
        return Response(resp.Error("ERR to few args"))
    }

    db := server.GetDB()
    key := args[0]
    
	logtitude, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return Response(resp.Error(err.Error()))
	}

	latitude, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return Response(resp.Error(err.Error()))
	}
	

	distance, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return Response(resp.Error(err.Error()))
	}

	entry, ok := db.Get(key)
	if !ok {
		return Response(resp.Error("Key doesnt exist"))
	}

	zset, ok := entry.AsZset()
	if !ok {
		return Response(resp.Error("Data doesnt exist"))
	}

	var names []string
	for _, v := range zset.Set {
		coord := helpers.DecodeGeo(uint64(v.Score))

		dist := helpers.GeoDist(latitude, logtitude, coord.Latitude, coord.Longitude)
		if dist < float64(distance) {
			names = append(names, v.Value)
		}
	}

    return Response(resp.ArrayString(names...))
}