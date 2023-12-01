package ucfunc

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	uc "gokratos/api/uc/v1"
)

func newHttpClient() (uc.UserCenterHTTPClient, func(), error) {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("127.0.0.1:8050"),
	)
	if err != nil {
		return nil, nil, err
	}

	close := func() {
		conn.Close()
	}

	client := uc.NewUserCenterHTTPClient(conn)
	return client, close, nil
}

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

func loginTemplate(appid, appsecret, aucode string, connType string) (string, error) {
	var reply *uc.LoginReply

	switch connType {
	case "http":
		client, close, err := newHttpClient()
		if err != nil {
			return "", err
		}
		defer close()
		ctx := context.Background()
		reply, err = client.Login(ctx, &uc.LoginRequest{
			AccessType: "code",
			AppId:      appid,
			AppSecret:  appsecret,
			AuCode:     aucode,
		})
		if err != nil {
			return "", err
		}
	case "grpc":
		client, close, err := newGrpcClient()
		if err != nil {
			return "", err
		}
		defer close()
		ctx := context.Background()
		reply, err = client.Login(ctx, &uc.LoginRequest{
			AccessType: "code",
			AppId:      appid,
			AppSecret:  appsecret,
			AuCode:     aucode,
		})
		if err != nil {
			return "", err
		}
	default:
	}

	if len(reply.AccessToken) == 0 {
		return "", errors.New("AccessToken is null string")
	}

	return reply.AccessToken, nil
}
