package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type InfoCommand struct {
}

func (c InfoCommand) Execute(args []string, db *database.DB, client *clients.Client) resp.Value {
	key := args[0]

	switch strings.ToLower(key) {
	case "replication":
		return resp.Bulk("# Replication \n role:master")
	}

	return resp.Bulk("")
}
