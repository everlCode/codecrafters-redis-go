package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Dir string
	Dbfilename string
	Port string
	ReplicaOf string
	Appendonly string
	Appenddirname string
	Appendfilename string
	Appendfsync string
}

func New() *Config {
	dirName, err := os.Getwd()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	port := flag.String("port", "6379", "redis port")
	replicaOf := flag.String("replicaof", "", "repica parameter")
	dir := flag.String("dir", dirName, "dir")
	dbfilename := flag.String("dbfilename", "", "file name")
	appendonly := flag.String("appendonly", "no", "no")
	appenddirname := flag.String("appenddirname", "appendonlydir", "appendonlydir")
	appendfilename := flag.String("appendfilename", "appendonly.aof", "appendonly.aof")
	appendfsync := flag.String("appendfsync", "everysec", "everysec")
	flag.Parse()

	return &Config{
		Dir: *dir,
		Dbfilename: *dbfilename,
		Port: *port,
		ReplicaOf: *replicaOf,
		Appendonly: *appendonly,
		Appenddirname: *appenddirname,
		Appendfilename: *appendfilename,
		Appendfsync: *appendfsync,
	}
}

