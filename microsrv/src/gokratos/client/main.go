package main

func main() {
	opt := "nacos"

	switch opt {
	case "simple":
		simpleHttp()
		simpleRpc()
	case "nacos":
		nacosRpc()
	default:
	}
}
