package main

import (
	"net"

	"github.com/jrieck1991/mark/internal/metrics"
	"github.com/jrieck1991/mark/internal/pipe"

	"google.golang.org/grpc"
)

const addr = ":7777"

func main() {

	// serve metrics via http
	go func() {
		if err := metrics.Serve(":9000"); err != nil {
			panic(err)
		}
	}()

	// listen on socket
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	p := pipe.Server{}
	s := grpc.NewServer()

	// serve requests
	pipe.RegisterPipeServer(s, &p)
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
