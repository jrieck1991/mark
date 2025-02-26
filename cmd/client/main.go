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
	msgSent   string = "message_sent"
	bytesSent string = "bytes_sent"
	sendErr   string = "send_error"
)

func main() {

	// serve metrics via http
	go func() {
		if err := metrics.Serve(metricsAddr); err != nil {
			panic(err)
		}
	}()

	// get metric counters
	counters := metrics.Counters(namespace, subsystem, []string{msgSent, bytesSent, sendErr})

	// dial server
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	pipeClient := pipe.NewPipeClient(conn)

	// get Ingest client
	client, err := pipeClient.Ingest(context.Background())
	if err != nil {
		panic(err)
	}
	defer client.CloseSend()

	fmt.Println("init complete")

	// send data to server
	for {

		d := &pipe.Data{
			Data: make([]byte, 1000),
		}

		if err := client.Send(d); err != nil {

			// counter send_error
			counters[sendErr].Inc()

			fmt.Println(err)
			continue
		}

		// counter message_sent
		counters[msgSent].Inc()

		// counter bytes_sent
		counters[bytesSent].Add(float64(len(d.GetData())))
	}
}
