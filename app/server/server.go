package server

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	MASTER_ROLE = "master"
	SLAVE_ROLE = "slave"
)

type Server struct {
	db *database.DB
	Pubsub *pubsub.Pubsub
	config *config.Config
	Port string
	Role string
	Replicas []*Replica
	MasterHost string
	MasterPort string
	MasterReplyId string
	MasterConnection net.Conn
	masterOffset int
}
func New() *Server {
	config := config.New()
	db := database.New(config)
	pubsub := pubsub.New()
	
	replicaData := strings.Split(config.ReplicaOf, " ")
	var masterHost string
	var MasterPort string
	if len(replicaData) > 1 {
		masterHost = replicaData[0]
		MasterPort = replicaData[1]
	}
	
	role := MASTER_ROLE
	if config.ReplicaOf != "" {
		role = SLAVE_ROLE
	}

	return &Server{
		db: db,
		Pubsub: pubsub,
		config: config,
		Port: config.Port,
		Role: role,
		MasterHost: masterHost,
		MasterPort: MasterPort,
		MasterReplyId: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
	}
}

func (s *Server) Start() net.Listener {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:" + s.Port)
	if err != nil {
		fmt.Println("Failed to bind to port " + s.Port)
		fmt.Println(err)
		os.Exit(1)
	}

	return listener
}

func (s *Server) InitReplica() *clients.Client {
	conn, err := net.Dial("tcp", s.MasterHost + ":" + s.MasterPort)
	if err != nil {
		panic(err)
	}
	s.MasterConnection = conn
	
	client := clients.New(conn)
	client.SetMasterConnection(true)

	s.MasterConnection = conn

	parser := client.GetParser()

	s.writeCommand(conn, "PING")
	parser.Read()

	s.writeCommand(
		conn,
		"REPLCONF",
		"listening-port",
		s.Port,
	)
	s.mustReadOk(parser)

	s.writeCommand(
		conn,
		"REPLCONF",
		"capa",
		"psync2",
	)
	s.mustReadOk(parser)

	s.writeCommand(
		conn,
		"PSYNC",
		"?",
		"-1",
	)

	parser.Read()

	parser.ReadRDB()
	client.SetOffset(0)

	return client
}

func (s *Server) AddOffset(v int) {
	s.masterOffset += v
}

func (s *Server) GetOffset() int {
	return s.masterOffset
}

func (s *Server) GetReplicas() []*Replica {
	return s.Replicas
}

func (s *Server) FindReplicaByClient(client *clients.Client) *Replica {
	for _, r := range s.GetReplicas() {
		if r.GetClient() == client {
			return  r
		}
	}

	return nil
}

func (s *Server) AddReplica(client *clients.Client) {
	replica := Replica{
		client: client,
	}

	s.Replicas = append(s.Replicas, &replica)
}

func (s *Server) writeCommand(
	conn net.Conn,
	arguments ...string,
) {
	values := make([]any, len(arguments))

	for i, arg := range arguments {
		values[i] = arg
	}

	_, err := conn.Write(
		resp.Array(values).Marshal(),
	)
	if err != nil {
		panic(err)
	}
}

func (s *Server) mustReadOk(parser *resp.Parser) {
	response, err := parser.Read()
	if err != nil {
		panic(err)
	}

	if !response.IsOk() {
		panic("invalid response")
	}
}

func (s *Server) SendRequest(client *clients.Client, arguments ...string) {
	conn := client.GetConnection()
	values := make([]any, len(arguments))

	for i, arg := range arguments {
		values[i] = arg
	}

	_, err := conn.Write(resp.Array(values).Marshal())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}

func (s *Server) SendRdb(conn net.Conn) {
	base64RDB := "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

	decoded, err := base64.StdEncoding.DecodeString(base64RDB)
	if err != nil {
		panic(err)
	}

	header := fmt.Sprintf("$%d\r\n", len(decoded))

	_, err = conn.Write([]byte(header))
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(decoded)
	if err != nil {
		panic(err)
	}
}

func (s *Server) SendPropagation(
	commandName string,
	args []string,
	server *Server,
) {
	request := append(
		[]string{commandName},
		args...,
	)

	data := resp.ArrayString(request).Marshal()
	server.AddOffset(len(data))

	for _, replica := range s.Replicas {
		_, err := replica.GetClient().GetConnection().Write(data)
		if err != nil {
			continue
		}
	}
}

func (s *Server) GetDB() *database.DB {
	return s.db
}

func (s *Server) GetConfig() *config.Config {
	return s.config
}

