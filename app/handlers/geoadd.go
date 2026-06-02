package handlers

import (
	"errors"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type GeoaddCommand struct {
}

func (c GeoaddCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 3 {
		return Response(resp.Error("Too few args!"))
	}
	var longtitude, latitude float64
	//key := args[0]
	longtitudeString := args[1]
	latitudeString := args[2]
	//member := args[3]

	longtitude, err := strconv.ParseFloat(longtitudeString, 64)
	if err != nil {
		return Response(resp.Error("Incorrect arg!"))
	}

	latitude, err = strconv.ParseFloat(latitudeString, 64)
	if err != nil {
		return Response(resp.Error("Incorrect arg!"))
	}

	err = validateCoords(longtitude, latitude)
	if err != nil {
		return Response(resp.Error(err.Error()))
	}
	
	return Response(resp.Integer(1))
}

func validateCoords(longtitude float64, latitude float64) (error) {
	var isErr bool
	if longtitude < -180 || longtitude > 180 {
		isErr = true
	}

	if latitude < -85.05112878 || latitude > 85.05112878 {
		isErr = true
	}

	if isErr {
		return errors.New("ERR invalid longitude,latitude pair 180.000000,90.000000")
	}

	return nil
}
