package server

import (
	"flag"

	"github.com/codecrafters-io/redis-starter-go/app/database"
)

type Server struct {
	db *database.DB
	config *Config
	Port string
	Role string
	MasterReplyId string
	MasterReplyOffset string
}

func New(db *database.DB) *Server {
	config := NewConfig()
	port := flag.String("port", "6379", "redis port")
	replicaOf := flag.String("replicaof", "", "repica parameter")
	flag.Parse()

	role := "master"
	if *replicaOf != "" {
		role = "slave"
	}
	return &Server{
		db: db,
		config: config,
		Port: *port,
		Role: role,
		MasterReplyId: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		MasterReplyOffset: "0",
	}
}

func (s *Server) GetDB() *database.DB {
	return s.db
}

func (s *Server) GetConfig() *Config {
	return s.config
}

