package main

import (
	"context"
	v1 "gokratos/api/greeter/v1"
	"log"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

func simpleHttp() {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "http yes!"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}

func simpleRpc() {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := v1.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "grpc yes!"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[grpc] SayHello %+v\n", reply.Message)
}
