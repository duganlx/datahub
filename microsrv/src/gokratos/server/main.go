package main

import (
	"context"
	"fmt"
	"log"

	helloworld "gokratos/helloworld/v1"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
)

var (
	Name    = "srv1"
	Version = "v1.0.0"

	NacosIp = "192.168.0.104"
	NacosNamespaceId = "2fe77cf0-7920-4405-82e6-bea518447a2f"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("server panic")
	}
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func nestedCnf(httpaddr, grpcaddr string) {
	s := &server{}
	httpSrv := http.NewServer(
		http.Address(httpaddr),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(grpcaddr),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)
	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	// nacos
	// == begin ==
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(NacosIp, 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId: NacosNamespaceId,
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  cc,
		},
	)

	if err != nil {
		panic(err)
	}
	r := nacos.New(client,
		nacos.WithGroup("groupA"),
		nacos.WithCluster("clusterA"),
	)
	// == end ==

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func localCnf(path string) (string, string) {
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var v struct {
		Server struct {
			Http struct {
				Addr    string `json:"addr"`
				Timeout string `json:"Timeout"`
			} `json:"http"`
			Grpc struct {
				Addr    string `json:"addr"`
				Timeout string `json:"Timeout"`
			} `json:"grpc"`
		} `json:"Server"`
	}

	if err := c.Scan(&v); err != nil {
		panic(err)
	}

	fmt.Printf("读取的配置为: %+v\n", v)
	return v.Server.Http.Addr, v.Server.Grpc.Addr
}

func nacosCnf() (string, string) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(NacosIp, 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId: NacosNamespaceId,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  cc,
		},
	)

	if err != nil {
		panic(err)
	}

	c := config.New(
		config.WithSource(
			knacos.NewConfigSource(
				client,
				knacos.WithGroup("groupA"),
				knacos.WithDataID("srv1conf.yaml"),
			),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}
	httpaddr, err := c.Value("server.http.addr").String()
	if err != nil {
		panic(err)
	}
	grpcaddr, err := c.Value("server.grpc.addr").String()
	if err != nil {
		panic(err)
	}

	return httpaddr, grpcaddr
}

func main() {

	opt := "nacos"

	switch opt {
	case "nested":
		nestedCnf(":8000", ":9000")
	case "local":
		httpaddr, grpcaddr := localCnf("./config.yaml")
		nestedCnf(httpaddr, grpcaddr)
	case "nacos":
		httpaddr, grpcaddr := nacosCnf()
		nestedCnf(httpaddr, grpcaddr)
	default:
	}

}
