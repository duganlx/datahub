package main

import (
	"gokratos/client/porder"
	"gokratos/client/ucfunc"
)

func main() {
	opt := "placeorder"

	switch opt {
	case "authfunc":
		// 测试鉴权功能
		ucfunc.Demo1("http")
		ucfunc.Demo2("http")
		ucfunc.Demo3("http")
		ucfunc.Demo4("http")
		ucfunc.Demo5("http")
		ucfunc.Demo6("http")
		ucfunc.Demo7("http")
		ucfunc.Demo8("http")
		ucfunc.Demo9("http")
		ucfunc.Demo10("http")

		ucfunc.Demo1("grpc")
		ucfunc.Demo2("grpc")
		ucfunc.Demo3("grpc")
		ucfunc.Demo4("grpc")
		ucfunc.Demo5("grpc")
		ucfunc.Demo6("grpc")
		ucfunc.Demo7("grpc")
		ucfunc.Demo8("grpc")
		ucfunc.Demo9("grpc")
		ucfunc.Demo10("grpc")
	case "placeorder":
		porder.PlaceOrder1("http")
		porder.PlaceOrder1("grpc")
	case "nacos":
		// 测试
		nacosRpc()
	default:
	}
}
