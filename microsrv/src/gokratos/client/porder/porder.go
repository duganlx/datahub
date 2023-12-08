package porder

import (
	"context"
	"errors"
	"fmt"
	au "gokratos/api/au/v1"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

func newAu1HttpClient() (au.AssetUnitHTTPClient, func(), error) {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		return nil, nil, err
	}

	close := func() {
		conn.Close()
	}

	client := au.NewAssetUnitHTTPClient(conn)
	return client, close, nil
}

func newAu1GrpcClient() (au.AssetUnitClient, func(), error) {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
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

	client := au.NewAssetUnitClient(conn)
	return client, close, nil
}

func send(appid, appsecret,  connType string) (bool, error) {
	var reply *au.PlaceOrderReply
	ctx := context.Background()

	switch connType {
	case "http":
		client, close, err := newAu1HttpClient()
		if err != nil {
			return false, err
		}
		defer close()
		reply, err = client.PlaceOrder(ctx, &au.PlaceOrderRequest{
			AppId:     appid,
			AppSecret: appsecret,
			OrderMsg:  "贵州茅台买100手",
		})
		if err != nil {
			return false, err
		}
	case "grpc":
		client, close, err := newAu1GrpcClient()
		if err != nil {
			return false, err
		}
		defer close()
		reply, err = client.PlaceOrder(ctx, &au.PlaceOrderRequest{
			AppId:     appid,
			AppSecret: appsecret,
			OrderMsg:  "贵州茅台买100手",
		})
		if err != nil {
			return false, err
		}
	default:
	}

	if !reply.Ok {
		return false, errors.New("place order fail")
	}

	return true, nil
}

// PlaceOrder1 用户模型要去资产单元下单整体过程
func PlaceOrder1(conntype string) {
	var appid = "kfuks"
	var appsecret = "4fd1ufklnksbry9"

	ok, err := send(appid, appsecret, conntype)
	if !ok || err != nil {
		fmt.Printf("[%s] %v\n", conntype, err)
		return
	}

	fmt.Printf("[%s] Place Order Success\n", conntype)
}
