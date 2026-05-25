package server

import "net"

type Replica struct {
	Connection net.Conn
}