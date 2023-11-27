package main

import (
	"context"
	"fmt"
	v1 "gokratos/helloworld/v1"
	"log"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	NacosIp          = "192.168.15.42"
	NacosPort        = uint64(8848)
	NacosNamespaceId = "ad13b472-04a0-4cf5-a4ee-d8cfd4cf81f5"
)

func nacosRpc() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(NacosIp, NacosPort),
	}

	cc := constant.ClientConfig{
		NamespaceId:          NacosNamespaceId,
		TimeoutMs:            5000,
		AppName:              "client1",
		OpenKMS:              false,
		LogDir:               "./log",
		LogLevel:             "debug",
		CacheDir:             "./cache",
		NotLoadCacheAtStart:  false,
		UpdateCacheWhenEmpty: true,
	}

	inamingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  &cc,
		},
	)
	if err != nil {
		panic(err)
	}

	registry := nacos.New(inamingClient, nacos.WithGroup("groupA"), nacos.WithCluster("clusterA"))

	intances, err := registry.GetService(context.Background(), "srv1.grpc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("intances: %+v\n", intances)

	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("discovery:///srv1.grpc"),
		transgrpc.WithDiscovery(registry),
	)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := v1.NewGreeterClient(conn)

	cnt := 1
	for {
		reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: fmt.Sprintf("grpc yes! - %d", cnt)})
		if err != nil {
			panic(err)
		}

		log.Printf("[grpc] SayHello %+v\n", reply.Message)
		cnt += 1
		time.Sleep(time.Second)
	}
}
