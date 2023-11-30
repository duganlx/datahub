package main

func main() {
	opt := "unittest"

	switch opt {
	case "simple":
		simpleHttp()
		simpleRpc()
	case "nacos":
		nacosRpc()
	case "unittest":
		Login1()
	default:
	}
}
