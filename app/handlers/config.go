package handlers

import (
	"reflect"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type ConfigCommand struct {
}

func (c ConfigCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	if len(args) < 2 {
		return Response(resp.Bulk("OK"))
	}
	param := args[1]
	paramNmae := strings.Title(param)
	reflection := reflect.ValueOf(*server.GetConfig())
	field := reflection.FieldByName(paramNmae)
	if field.IsZero() {
		return Response(resp.Error("Field doesnt exist!"))
	}
	v := field.Interface()
	value := v.(string)
	

	return Response(resp.ArrayString(param, value))
}
