package server

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	MASTER_ROLE = "master"
	SLAVE_ROLE = "slave"
)

type Server struct {
	db *database.DB
	config *Config
	Port string
	Role string
	Replicas []Replica
	MasterHost string
	MasterPort string
	MasterReplyId string
	MasterReplyOffset string
	MasterConnection net.Conn
}

func New(db *database.DB) *Server {
	config := NewConfig()
	port := flag.String("port", "6379", "redis port")
	replicaOf := flag.String("replicaof", "", "repica parameter")
	flag.Parse()

	replicaData := strings.Split(*replicaOf, " ")
	var masterHost string
	var MasterPort string
	if len(replicaData) > 1 {
		masterHost = replicaData[0]
		MasterPort = replicaData[1]
	}
	
	role := MASTER_ROLE
	if *replicaOf != "" {
		role = SLAVE_ROLE
	}
	return &Server{
		db: db,
		config: config,
		Port: *port,
		Role: role,
		MasterHost: masterHost,
		MasterPort: MasterPort,
		MasterReplyId: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplyOffset: "0",
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

	parser := resp.New(conn)

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

	fullResync, err := parser.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(fullResync.String)

	s.readRdb(parser)

	return client
}

func (s *Server) AddReplica(conn net.Conn) {
	replica := Replica{
		Connection: conn,
	}

	s.Replicas = append(s.Replicas, replica)
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

func (s *Server) readRdb(parser *resp.Parser) {
	parser.Read()
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

func (s *Server) SendRequest(conn net.Conn, arguments ...string) resp.Value {
	parser := resp.New(conn)
	values := make([]any, len(arguments))

	for i, arg := range arguments {
		values[i] = arg
	}

	_, err := conn.Write(resp.Array(values).Marshal())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	response, err := parser.Read()
	
	return response
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
) {
	request := append(
		[]string{commandName},
		args...,
	)

	data := resp.ArrayString(request).Marshal()

	for _, replica := range s.Replicas {
		_, err := replica.Connection.Write(data)
		if err != nil {
			continue
		}
	}
}

func (s *Server) GetDB() *database.DB {
	return s.db
}

func (s *Server) GetConfig() *Config {
	return s.config
}

