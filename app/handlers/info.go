package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type InfoCommand struct {
}

func (c InfoCommand) Execute(
	args []string,
	server *server.Server,
	client *clients.Client,
) CommandResponse {

	//section := strings.ToLower(args[0])

	info :=
		"# Replication\r\n" +
			"role:" + server.Role + "\r\n" +
			"master_replid:" + server.MasterReplyId + "\r\n" +
			"master_repl_offset:" + "0" + "\r\n"

	return Response(resp.Bulk(info))

}
