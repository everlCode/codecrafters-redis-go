package server

import (
	"github.com/codecrafters-io/redis-starter-go/app/database"
)

type Server struct {
	db *database.DB
	config *Config
}

func New(db *database.DB) *Server {
	config := NewConfig()
	return &Server{
		db: db,
		config: config,
	}
}

func (s *Server) GetDB() *database.DB {
	return s.db
}

func (s *Server) GetConfig() *Config {
	return s.config
}

