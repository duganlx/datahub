package main

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	uc "gokratos/api/uc/v1"
)

// 前提说明：
// - 用户级的访问令牌

func newClient() (uc.UserCenterHTTPClient, func(), error) {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("127.0.0.1:8050"),
	)
	close := func() {
		conn.Close()
	}

	if err != nil {
		return nil, nil, err
	}
	
	client := uc.NewUserCenterHTTPClient(conn)
	return client, close, nil
}

// Login1 登录场景1
// 用户访问
func Login1() {
	client, close, err := newClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer close()
	ctx := context.Background()

	reply, err := client.Login(ctx, &uc.LoginRequest{
		AccessType: "code",
		AppId:      "asdj",
		AppSecret:  "d54sdfejbd561sa",
		AuCode:     "0148P1016_ww",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(reply)
}
