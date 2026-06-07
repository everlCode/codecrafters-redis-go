package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GeodistCommand struct {
}

func (c GeodistCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 3 {
		return Response(resp.Error("to few args"))
	}
	db := server.GetDB()

	key := args[0]
	member1 := args[1]
	member2 := args[2]

	entry, ok := db.Get(key)
	if !ok {
		return Response(resp.Error("dfddf"))
	}
	value, ok := entry.AsZset()
	if !ok {
		return Response(resp.Error("dfddf"))
	}
	firstMember := value.Get(member1)
	secondMember := value.Get(member2)

	if firstMember == nil || secondMember == nil {
		return Response(resp.Error("fgfg"))
	}

	coord1 := helpers.DecodeGeo(uint64(firstMember.Score))
	coord2 := helpers.DecodeGeo(uint64(secondMember.Score))

	res := helpers.GeoDist(coord1.Latitude, coord1.Longitude, coord2.Latitude, coord2.Longitude)
	

	return Response(resp.Bulk(strconv.FormatFloat(res, 'f', 4, 64)))
}
