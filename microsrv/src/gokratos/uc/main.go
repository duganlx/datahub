package main

import (
	"context"
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

func unitTest(s *server.Server) {

	ctx := context.Background()

	// 用户: ww(15739) xjw(15743) wsy(15747) yrl(15753)
	// 部门: admin(10), operations(11), it(12)
	testSet := []struct {
		AuCode string
		OpType string
		Uid    int64
		Expect bool // 预期结果
	}{
		// 用户 ww 访问 0148P1016_ww 操作类型 r  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "r", Uid: 15739, Expect: true},
		// 用户 ww 访问 0148P1016_ww 操作类型 w  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "w", Uid: 15739, Expect: true},
		// 用户 ww 访问 0148P1016_ww 操作类型 x  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "x", Uid: 15739, Expect: true},
		// 用户 yrl 访问 0148P1016_ww 操作类型 r  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "r", Uid: 15753, Expect: true},
		// 用户 yrl 访问 0148P1016_ww 操作类型 w  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "w", Uid: 15753, Expect: true},
		// 用户 yrl 访问 0148P1016_ww 操作类型 x  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "x", Uid: 15753, Expect: true},
		// 用户 ww 访问 88853899_ww 操作类型 r  ==>  成功
		{AuCode: "88853899_ww", OpType: "r", Uid: 15739, Expect: true},
		// 用户 ww 访问 88853899_ww 操作类型 w  ==>  失败，对该资产单元无该操作类型权限
		{AuCode: "88853899_ww", OpType: "w", Uid: 15739, Expect: false},
		// 用户 ww 访问 121000 操作类型 w  ==>  失败，对该资产单元无任何访问权限
		{AuCode: "121000", OpType: "w", Uid: 15739, Expect: false},
		// 用户 wsy 访问 DRW001ZTX_04 操作类型 x  ==>  成功
		{AuCode: "DRW001ZTX_04", OpType: "x", Uid: 15747, Expect: true},
		// 用户 xjw 访问 EAMLS1ZT_00 操作类型 r  ==>  成功
		{AuCode: "EAMLS1ZT_00", OpType: "r", Uid: 15743, Expect: true},
		// 用户 xjw 访问 EAMLS1ZT_00 操作类型 w  ==>  成功
		{AuCode: "EAMLS1ZT_00", OpType: "w", Uid: 15743, Expect: true},
		// 用户 xjw 访问 EAMLS1ZT_00 操作类型 x  ==>  失败，对该资产单元仅有rw操作权限
		{AuCode: "EAMLS1ZT_00", OpType: "x", Uid: 15743, Expect: false},
		// 用户 xjw 访问 EAMLS1ZT_00 操作类型 所有*  ==>  失败，对该资产单元仅有rw操作权限
		{AuCode: "EAMLS1ZT_00", OpType: "*", Uid: 15743, Expect: false},
		// 部门 admin 访问 EAMLS1ZT_00 操作类型 x  ==>  成功
		{AuCode: "EAMLS1ZT_00", OpType: "x", Uid: 10, Expect: true},
		// 部门 operations 访问 0148P1016_ww 操作类型 r  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "r", Uid: 11, Expect: true},
		// 部门 operations 访问 0148P1016_ww 操作类型 w  ==>  失败，部门对所有资产单元只有r操作权限
		{AuCode: "0148P1016_ww", OpType: "w", Uid: 11, Expect: false},
		// 部门 it 访问 0148P1016_ww 操作类型 r  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "r", Uid: 12, Expect: true},
		// 部门 it 访问 0148P1016_ww 操作类型 w  ==>  成功
		{AuCode: "0148P1016_ww", OpType: "w", Uid: 12, Expect: true},
		// 部门 it 访问 121000 操作类型 r  ==>  成功
		{AuCode: "121000", OpType: "r", Uid: 12, Expect: true},
		// 部门 it 访问 121000 操作类型 w  ==>  失败，部门对该资产单元仅有r操作权限
		{AuCode: "121000", OpType: "w", Uid: 12, Expect: false},
	}

	for index, unit := range testSet {
		in := &uc.AuthrawRequest{
			AuCode: unit.AuCode,
			OpType: unit.OpType,
			Uid:    unit.Uid,
		}
		reply, err := s.Authraw(ctx, in)
		if err != nil {
			fmt.Printf("[Demo%d] X: %+v\n", index, err)
			continue
		}
		if unit.Expect != reply.Ok {
			fmt.Printf("[Demo%d] X: %+v\n", index, unit)
			continue
		}

		fmt.Printf("[Demo%d] √\n", index)
	}

}

func main() {

	opt := "unitTest"

	switch opt {
	case "unitTest":
		s, err := server.NewServer()
		if err != nil {
			panic(err)
		}

		unitTest(s)
		s.Cbe.UpdateAuth()
		unitTest(s)

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
