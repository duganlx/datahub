package main

import (
	"log"

	au "gokratos/api/au/v1"
	"gokratos/au1/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	// "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	// "github.com/nacos-group/nacos-sdk-go/clients"
	// "github.com/nacos-group/nacos-sdk-go/common/constant"
	// "github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	Name    = "au1"
	Version = "v1.0.0"

	NacosIp          = "192.168.15.42"
	NacosNamespaceId = "ad13b472-04a0-4cf5-a4ee-d8cfd4cf81f5"
)

func startSrv(httpaddr, grpcaddr string) {
	s := server.NewAuServer()
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
	au.RegisterAssetUnitServer(grpcSrv, s)
	au.RegisterAssetUnitHTTPServer(httpSrv, s)

	// nacos
	// == begin ==
	// sc := []constant.ServerConfig{
	// 	*constant.NewServerConfig(NacosIp, 8848),
	// }

	// cc := &constant.ClientConfig{
	// 	NamespaceId: NacosNamespaceId,
	// }

	// client, err := clients.NewNamingClient(
	// 	vo.NacosClientParam{
	// 		ServerConfigs: sc,
	// 		ClientConfig:  cc,
	// 	},
	// )

	// if err != nil {
	// 	panic(err)
	// }
	// r := nacos.New(client,
	// 	nacos.WithGroup("groupA"),
	// 	nacos.WithCluster("clusterA"),
	// )
	// == end ==

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		// kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
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

// func nacosCnf() (string, string) {
// 	sc := []constant.ServerConfig{
// 		*constant.NewServerConfig(NacosIp, 8848),
// 	}

// 	cc := &constant.ClientConfig{
// 		NamespaceId: NacosNamespaceId,
// 	}

// 	client, err := clients.NewConfigClient(
// 		vo.NacosClientParam{
// 			ServerConfigs: sc,
// 			ClientConfig:  cc,
// 		},
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	c := config.New(
// 		config.WithSource(
// 			knacos.NewConfigSource(
// 				client,
// 				knacos.WithGroup("groupA"),
// 				knacos.WithDataID("srv1conf.yaml"),
// 			),
// 		),
// 	)

// 	if err := c.Load(); err != nil {
// 		panic(err)
// 	}
// 	httpaddr, err := c.Value("server.http.addr").String()
// 	if err != nil {
// 		panic(err)
// 	}
// 	grpcaddr, err := c.Value("server.grpc.addr").String()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return httpaddr, grpcaddr
// }

func main() {

	opt := "nested"

	switch opt {
	case "nested":
		startSrv(":8000", ":9000")
	// case "local":
	// 	httpaddr, grpcaddr := localCnf("./config.yaml")
	// 	nestedCnf(httpaddr, grpcaddr)
	// case "nacos":
	// 	httpaddr, grpcaddr := nacosCnf()
	// 	nestedCnf(httpaddr, grpcaddr)
	default:
	}

}
