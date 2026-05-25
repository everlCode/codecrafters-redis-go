package server

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

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
	MasterHost string
	MasterPort string
	MasterReplyId string
	MasterReplyOffset string
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

	if s.Role == SLAVE_ROLE {
		s.initReplica()
	}

	return listener
}

func (s *Server) initReplica() {
	conn, err := net.Dial("tcp", s.MasterHost + ":" + s.MasterPort)
	if err != nil {
		panic(err)
	}

	s.sendRequest(conn, "PING")
	
	firstRequest := s.sendRequest(conn, "REPLCONF", "listening-port", s.Port)
	if !firstRequest.IsOk() {
		panic(errors.New("Replica start error"))
	}
	
	secondRequest := s.sendRequest(conn, "REPLCONF", "capa", "psync2")
	if !secondRequest.IsOk() {
		panic(errors.New("Replica start error"))
	}

	s.sendRequest(conn, "PSYNC", "?", "-1")
}

func (s *Server) sendRequest(conn net.Conn, arguments ...string) resp.Value {
	parser := resp.New(conn)
	values := make([]any, len(arguments))

	for i, arg := range arguments {
		values[i] = arg
	}

	_, err := conn.Write(resp.Array(values).Marshal())
	if err != nil {
		panic(err)
	}

	response, err := parser.Read()
	
	return response
}

func (s *Server) GetDB() *database.DB {
	return s.db
}

func (s *Server) GetConfig() *Config {
	return s.config
}

