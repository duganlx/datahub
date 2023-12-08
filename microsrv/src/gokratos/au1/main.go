package main

import (
	"fmt"
	"log"

	au "gokratos/api/au/v1"
	"gokratos/au1/ncs"
	"gokratos/au1/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
)

var (
	Name    = "au1"
	Version = "v1.0.0"
)

type Cfg struct {
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

	Au struct {
		Code string `json:"code"`
	} `json:"au"`
}

func startSrv(cfgsrc string) {
	var httpaddr, grpcaddr, aucode string

	switch cfgsrc {
	case "nested":
		httpaddr = ":8000"
		grpcaddr = ":9000"
		aucode = "EAMLS1ZT_00"
	case "local":
		path := "./config.yaml"
		c := config.New(
			config.WithSource(
				file.NewSource(path),
			),
		)

		if err := c.Load(); err != nil {
			panic(err)
		}

		var v Cfg

		if err := c.Scan(&v); err != nil {
			panic(err)
		}

		fmt.Printf("读取的配置为: %+v\n", v)
		httpaddr = v.Server.Http.Addr
		grpcaddr = v.Server.Grpc.Addr
		aucode = v.Au.Code
	case "nacos":
		cfgcli, err := ncs.NewNacosConfigClient()
		if err != nil {
			panic(err)
		}

		c := config.New(
			config.WithSource(
				knacos.NewConfigSource(
					cfgcli,
					knacos.WithGroup("groupA"),
					knacos.WithDataID("srv1conf.yaml"),
				),
			),
		)

		if err := c.Load(); err != nil {
			panic(err)
		}
		httpaddr, err = c.Value("server.http.addr").String()
		if err != nil {
			panic(err)
		}
		grpcaddr, err = c.Value("server.grpc.addr").String()
		if err != nil {
			panic(err)
		}
		aucode, err = c.Value("au.code").String()
		if err != nil {
			panic(err)
		}
	default:
	}

	s := server.NewAuServer(aucode)
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

	incli, err := ncs.NewNacosInamingClient()
	if err != nil {
		panic(err)
	}
	r := nacos.New(incli,
		nacos.WithGroup("groupA"),
		nacos.WithCluster("clusterA"),
	)

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
	// nested local nacos
	startSrv("nacos")
}
