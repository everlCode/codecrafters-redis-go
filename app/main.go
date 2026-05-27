package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/dispatcher"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	db := database.New()
	s := server.New(db)
	listener := s.Start()

	var masterClient *clients.Client
	if s.Role == server.SLAVE_ROLE {
		masterClient = s.InitReplica()
	}

	if s.MasterConnection != nil {
		go handle(masterClient, s)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		client := clients.New(conn)
		go handle(client, s)
	}
}

func handle(client *clients.Client, server *server.Server) {
	conn := client.GetConnection()
	defer conn.Close()

	writer := NewWriter(conn)
	parser := client.GetParser()

	for {
		request, err := parser.Read()
		
		if err != nil {
			if err == io.EOF {
				return
			}
			writer.Write(resp.Error(
				fmt.Sprintf("ERR: %s", err.Error()),
			))
			continue
		}
		if request.Type != resp.ARRAY {
			writer.Write(resp.Error(
				fmt.Sprintf("ERR: %s", "Неверный запрос! Ожидался массив!"),
			))
			continue
		}
		if len(request.Array) == 0 {
			writer.Write(resp.Error(
				fmt.Sprintf("ERR: %s", "Неверный запрос! Не переданы необходимые аргументы"),
			))

			continue
		}
		args := resp.ParseSlice(request.Array)
		result := dispatcher.Dispatch(args, server, client)

		if (!client.MasterConnection() && !client.IsReplica()) || result.NeedAnswer {
			writer.Write(result.GetResponse())
		}

		if client.IsReplica() && !client.RDBSent {
			server.SendRdb(client.GetConnection())
			client.RDBSent = true

			server.AddReplica(client)

			// server.SendRequest(client.GetConnection(), "REPLCONF", "GETACK", "*")
			// server.SendRequest(client.GetConnection(), "SET", "strawberry", "pear")
			// server.SendRequest(client.GetConnection(), "SET", "orange", "apple")
			// server.SendRequest(client.GetConnection(), "REPLCONF", "GETACK", "*")
		}

		if result.IsPropagation && !client.MasterConnection() {
			server.SendPropagation(result.PropCommand.Name, result.PropCommand.Args, server)
		}

		if client.MasterConnection() {
			client.AddOffset(parser.GetOffset())
		}
	}
}
