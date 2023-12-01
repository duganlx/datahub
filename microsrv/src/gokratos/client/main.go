package main

func main() {
	opt := "direct"

	switch opt {
	case "direct":
		// 测试鉴权功能
		Demo1("http")
		Demo2("http")
		Demo3("http")
		Demo4("http")
		Demo5("http")
		Demo6("http")
		Demo7("http")
		Demo8("http")
		Demo9("http")
		Demo10("http")

		Demo1("grpc")
		Demo2("grpc")
		Demo3("grpc")
		Demo4("grpc")
		Demo5("grpc")
		Demo6("grpc")
		Demo7("grpc")
		Demo8("grpc")
		Demo9("grpc")
		Demo10("grpc")
	case "nacos":
		// 测试
		nacosRpc()
	default:
	}
}
