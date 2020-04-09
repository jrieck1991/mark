package main

import (
	"context"

	"fmt"

	"github.com/jrieck1991/mark/internal/metrics"
	"github.com/jrieck1991/mark/internal/pipe"
	"google.golang.org/grpc"
)

const (
	serverAddr  string = "server:7777"
	metricsAddr string = ":9000"

	// metrics
	namespace string = "app"
	subsystem string = "client"
)

func main() {

	// serve metrics via http
	go func() {
		if err := metrics.Serve(metricsAddr); err != nil {
			panic(err)
		}
	}()

	counters := metrics.Counters(namespace, subsystem, []string{})

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
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
