package server

import (
	"github.com/codecrafters-io/redis-starter-go/app/clients"
)

type Replica struct {
	client *clients.Client
	offset int
}

func (r *Replica) GetClient() *clients.Client {
	return r.client
}

func (r *Replica) SetOffset(v int) {
	r.offset = v
}

func (r *Replica) GetOffset() int {
	return r.offset
}


