package handlers

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/clients"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

type WaitCommand struct {
}

func (c WaitCommand) Execute(args []string, server *server.Server, client *clients.Client) CommandResponse {
	var replicaCount, timeout int
	if len(args) > 0 {
		replicaCountStr := args[0]
		count, err := strconv.Atoi(replicaCountStr)
		replicaCount = count
		if err != nil {
			return Response(resp.Error(err.Error()))
		}
	}

	if len(args) > 1 {
		timeoutStr := args[1]
		t, err := strconv.Atoi(timeoutStr)
		timeout = t
		if err != nil {
			return Response(resp.Error(err.Error()))
		}
	}


	var count int = 0
	replicas := server.GetReplicas()

	if len(replicas) > 0 {
		replicas := server.GetReplicas()
		
		for i := 0; i < len(replicas); i++ {
			replic := replicas[i]
			server.SendRequest(replic.GetClient().GetConnection(), "REPLCONF", "GETACK", "*")
		}

		now := time.Now()
		deadline := now.Add(time.Duration(timeout) * time.Millisecond)
		targetOffset := server.GetOffset()

		for time.Now().Before(deadline) {
			count = 0

			for _, replic := range replicas {
				if replic.GetOffset() >= targetOffset {
					count++
				}
			}

			if count >= replicaCount {
				break
			} 
		}

	}

	return Response(resp.Integer(count))
}
