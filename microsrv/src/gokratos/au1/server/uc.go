package server

import "fmt"

func checkAuth(appid, appsecret, aucode string) (bool, error) {
	// todo 去用户中心鉴权
	fmt.Println("checkAuth...")

	return true, nil
}
