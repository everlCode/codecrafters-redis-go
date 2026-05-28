package server

import "flag"

type Config struct {
	Dir string
	Dbfilename string
	Port string
	ReplicaOf string
}

func NewConfig() *Config {
	port := flag.String("port", "6379", "redis port")
	replicaOf := flag.String("replicaof", "", "repica parameter")
	dir := flag.String("dir", "/tmp/redis-data", "dir")
	dbfilename := flag.String("dbfilename", "dump.rdb", "file name")
	flag.Parse()

	return &Config{
		Dir: *dir,
		Dbfilename: *dbfilename,
		Port: *port,
		ReplicaOf: *replicaOf,
	}
}

