package main

import (
	"net"
	"sandbox/grpc/pipe"

	"google.golang.org/grpc"
)

const addr = "localhost:7777"

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
