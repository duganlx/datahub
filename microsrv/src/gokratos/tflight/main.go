package main

import (
	"github.com/apache/arrow/go/v15/arrow/flight"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type flightServer struct {
	mem memory.Allocator
	flight.BaseFlightServer
}

func main() {
	s := flight.NewFlightServer()
	s.Init("localhost:0")
	f := &flightServer{}
	s.RegisterFlightService(f)

	go s.Serve()
	defer s.Shutdown()

	client, err := flight.NewFlightClient(s.Addr().String(), nil, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// info, err := fli
}
