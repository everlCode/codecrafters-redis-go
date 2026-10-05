package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/helpers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

const (
	AND = "AND"
)

type BitOpCommand struct {
}

func (c BitOpCommand) Execute(
	args []string,
	server *server.Server,
	client *clients.Client,
) CommandResponse {
	if len(args) < 4 {
		return Response(resp.Error("ERR to few args"))
	}
	db := server.GetDB()

	operation := args[0]
	curKey := args[1]
	firstKey := args[2]
	secondKey := args[3]

	var curEntry database.Entry
	curEntry, ok := db.Get(curKey)
	if !ok {
		curEntry = database.String("")
	}

	firstEntry, ok := db.Get(firstKey)
	if !ok {
		firstEntry = database.String("")
	}

	secondEntry, ok := db.Get(secondKey)
	if !ok {
		secondEntry = database.String("")
	}

	var result []byte
	switch strings.ToUpper(operation) {
	case AND:
		result = helpers.BitAndOp([]byte(firstEntry.AsString()), []byte(secondEntry.AsString()))
	}

	curEntry.Set(string(result))

	db.Set(curKey, curEntry)

	response := Response(resp.Integer(len(result)))
	//response.SetPropagationCommand(SETBIT, args)

	return response
}
