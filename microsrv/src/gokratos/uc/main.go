package main

import (
	"context"
	uc "gokratos/api/uc/v1"
	"log"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	Name    = "uc"
	Version = "v1.0.0"

	NacosIp          = "192.168.15.42"
	NacosNamespaceId = "ad13b472-04a0-4cf5-a4ee-d8cfd4cf81f5"
)

type server struct {
	uc.UnimplementedUserCenterServer
}

func (s *server) Login(ctx context.Context, in *uc.LoginRequest) (*uc.LoginReply, error) {

	return &uc.LoginReply{
		AccessToken:  "xxx",
		RefreshToken: "kkk",
		TokenType:    "ttt",
		Expires:      120,
		Scrope:       "",
		Uid:          100,
	}, nil
}

func startSrv(httpaddr, grpcaddr string) {
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
	uc.RegisterUserCenterServer(grpcSrv, s)
	uc.RegisterUserCenterHTTPServer(httpSrv, s)

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

func main() {

	opt := "nested"

	switch opt {
	case "nested":
		startSrv(":8050", ":9050")
	// case "local":
	// 	httpaddr, grpcaddr := localCnf("./config.yaml")
	// 	nestedCnf(httpaddr, grpcaddr)
	// case "nacos":
	// 	httpaddr, grpcaddr := nacosCnf()
	// 	nestedCnf(httpaddr, grpcaddr)
	default:
	}

}
