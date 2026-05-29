package config

import "flag"

type Config struct {
	Dir string
	Dbfilename string
	Port string
	ReplicaOf string
}

func New() *Config {
	port := flag.String("port", "6379", "redis port")
	replicaOf := flag.String("replicaof", "", "repica parameter")
	dir := flag.String("dir", "", "dir")
	dbfilename := flag.String("dbfilename", "", "file name")
	flag.Parse()

	return &Config{
		Dir: *dir,
		Dbfilename: *dbfilename,
		Port: *port,
		ReplicaOf: *replicaOf,
	}
}

