package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"sandbox/grpc/pipe"

	"google.golang.org/grpc"
)

const (
	addr        = "localhost:7777"
	payloadSize = 32
	numReqs     = 1000
)

func main() {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pipe.NewPipeClient(conn)

	ingestClient, err := c.Ingest(context.TODO())
	if err != nil {
		panic(err)
	}

	for i := 0; i < numReqs; i++ {

		payload := make([]byte, payloadSize)
		rand.Read(payload)

		if err := ingestClient.Send(&pipe.Data{Data: payload}); err != nil {
			panic(err)
		}
		fmt.Printf("client: sent payload %d sucess\n", i)
	}

}
