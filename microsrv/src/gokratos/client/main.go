package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "gokratos/helloworld/v1"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	NacosIp          = "192.168.0.104"
	NacosNamespaceId = "2fe77cf0-7920-4405-82e6-bea518447a2f"
)

func main() {
	opt := "nacos"

	switch opt {
	case "simple":
		simpleHttp()
		simpleRpc()
	case "nacos":
		nacosRpc()
	default:
	}

}

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

func nacosRpc() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(NacosIp, 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId: NacosNamespaceId,
		TimeoutMs:   5000,
		AppName:     "client1",
		// Endpoint:             "192.168.15.42:8848",
		OpenKMS:              false,
		LogDir:               "./log",
		LogLevel:             "debug",
		CacheDir:             "./cache",
		NotLoadCacheAtStart:  false,
		UpdateCacheWhenEmpty: true,
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  &cc,
		},
	)
	if err != nil {
		fmt.Println("--> 0", err)
		return
	}

	r := nacos.New(cli, nacos.WithGroup("groupA"), nacos.WithCluster("clusterA"))

	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("discovery:///srv1.grpc"),
		transgrpc.WithDiscovery(r),
	)

	if err != nil {
		fmt.Println("--> 1", err)
		return
	}
	defer conn.Close()

	client := v1.NewGreeterClient(conn)

	cnt := 1
	for {
		reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: fmt.Sprintf("grpc yes! - %d", cnt)})
		if err != nil {
			fmt.Println("--> 2", err)
			return
		}

		log.Printf("[grpc] SayHello %+v\n", reply.Message)
		cnt += 1
		time.Sleep(time.Second)
	}
}
