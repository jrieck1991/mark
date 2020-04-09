package main

import (
	"context"

	"fmt"

	"github.com/jrieck1991/mark/internal/pipe"
	"google.golang.org/grpc"
)

const addr = "server:7777"

func main() {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	pipeClient := pipe.NewPipeClient(conn)

	client, err := pipeClient.Ingest(context.Background())
	if err != nil {
		panic(err)
	}
	defer client.CloseSend()

	fmt.Println("init complete")
	for {

		d := &pipe.Data{
			Data: make([]byte, 1000),
		}

		if err := client.Send(d); err != nil {
			panic(err)
		}
	}
}
