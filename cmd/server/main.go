package main

import (
	"net"

	_ "github.com/jrieck1991/mark/internal/metrics"
	"github.com/jrieck1991/mark/internal/pipe"

	"google.golang.org/grpc"
)

const addr = ":7777"

func main() {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	p := pipe.Server{}
	s := grpc.NewServer()

	pipe.RegisterPipeServer(s, &p)
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
