package server

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	uc "gokratos/api/uc/v1"
)

func newGrpcClient() (uc.UserCenterClient, func(), error) {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9050"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	close := func() {
		conn.Close()
	}

	client := uc.NewUserCenterClient(conn)
	return client, close, nil
}

func checkAuth(ctx context.Context, appid, appsecret, aucode string) (bool, error) {
	client, close, err := newGrpcClient()
	if err != nil {
		return false, err
	}
	defer close()

	reply, err := client.Login(ctx, &uc.LoginRequest{
		AccessType: "code",
		AppId:      appid,
		AppSecret:  appsecret,
		AuCode:     aucode,
	}) 
	if err != nil {
		return false, err
	}

	if len(reply.AccessToken) == 0 {
		return false, errors.New("AccessToken is null string")
	}

	return true, nil
}
