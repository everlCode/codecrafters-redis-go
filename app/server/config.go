package server

import (
	"flag"
)

type Config struct {
	Port string
	Role string
}

func NewConfig() *Config {
	port := flag.String("port", "6379", "redis port")
	replicaOf := flag.String("replicaof", "", "repica parameter")
	flag.Parse()

	role := "master"
	if *replicaOf != "" {
		role = "slave"
	}

	return &Config{
		Port: *port,
		Role: role,
	}
}