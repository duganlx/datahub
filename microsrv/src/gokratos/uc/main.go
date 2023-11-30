package main

import (
	"fmt"
	uc "gokratos/api/uc/v1"
	"gokratos/uc/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// import (

// 	"log"

// 	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
// 	"github.com/go-kratos/kratos/v2/config"
// 	"github.com/go-kratos/kratos/v2/config/file"
// 	"github.com/nacos-group/nacos-sdk-go/clients"
// 	"github.com/nacos-group/nacos-sdk-go/common/constant"
// 	"github.com/nacos-group/nacos-sdk-go/vo"
// )

var (
	Name    = "uc_v1"
	Version = "v1.0.0"

	NacosIp          = "192.168.15.42"
	NacosNamespaceId = "ad13b472-04a0-4cf5-a4ee-d8cfd4cf81f5"

	// Authentication
	ClientId     = "thisisaclientid"
	ClientSecret = "thisisaclientsecret"
	AppId        = "thisisaappid"
	AppSecret    = "thisisaappsecret"

	// User Info => mock user db
	UserId     = 1427818295636267008
	UserName   = "Tom"
	UserMobile = "15306218464"
)

func startSrv(httpaddr, grpcaddr string) {
	s := &server.Server{}

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

	// 	// nacos
	// 	// == begin ==
	// 	sc := []constant.ServerConfig{
	// 		*constant.NewServerConfig(NacosIp, 8848),
	// 	}

	// 	cc := &constant.ClientConfig{
	// 		NamespaceId: NacosNamespaceId,
	// 	}

	// 	client, err := clients.NewNamingClient(
	// 		vo.NacosClientParam{
	// 			ServerConfigs: sc,
	// 			ClientConfig:  cc,
	// 		},
	// 	)

	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	r := nacos.New(client,
	// 		nacos.WithGroup("groupA"),
	// 		nacos.WithCluster("clusterA"),
	// 	)
	// 	// == end ==

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		// kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}

// func localCnf(path string) (string, string) {
// 	c := config.New(
// 		config.WithSource(
// 			file.NewSource(path),
// 		),
// 	)

// 	if err := c.Load(); err != nil {
// 		panic(err)
// 	}

// 	var v struct {
// 		Server struct {
// 			Http struct {
// 				Addr    string `json:"addr"`
// 				Timeout string `json:"Timeout"`
// 			} `json:"http"`
// 			Grpc struct {
// 				Addr    string `json:"addr"`
// 				Timeout string `json:"Timeout"`
// 			} `json:"grpc"`
// 		} `json:"Server"`
// 	}

// 	if err := c.Scan(&v); err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("读取的配置为: %+v\n", v)
// 	return v.Server.Http.Addr, v.Server.Grpc.Addr
// }

func main() {

	opt := "nested"

	switch opt {
	case "nested":
		startSrv(":8050", ":9050")
	// case "local":
	// 	httpaddr, grpcaddr := localCnf("./config.yaml")
	// 	startSrv(httpaddr, grpcaddr)
	// case "nacos":
	// 	httpaddr, grpcaddr := nacosCnf()
	// 	nestedCnf(httpaddr, grpcaddr)
	default:
	}

}
