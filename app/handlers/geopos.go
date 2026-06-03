package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GeoposCommand struct {
}

func (c GeoposCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
    if len(args) < 2 {
        return Response(resp.Error("ERR to few args"))
    }

    db := server.GetDB()
    key := args[0]
    members := args[1:]

    var response []any
    for _, member := range members {
        entry, ok := db.Get(key)
        if !ok {
            // Ключ не существует — nil для этого члена
            response = append(response, resp.NullArray())
            continue
        }

        zset, _ := entry.AsZset()
        value := zset.Get(member)
        if value == nil {
            response = append(response, resp.NullArray())
            continue
        }

        //lon, lat := geo.DecodeScore(score)
        response = append(response, []any{
            0.0,
            0.0,
        })
    }

    return Response(resp.Array(response))
}