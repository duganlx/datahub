package server

import (
	"context"
	"fmt"
	au "gokratos/api/au/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

type server struct {
	au.UnimplementedAssetUnitServer

	akpair map[string]string // appsecret: appid
}

func NewAuServer() *server {
	s := &server{
		akpair: make(map[string]string),
	}

	return s
}

func (s *server) SayHello(ctx context.Context, in *au.HelloRequest) (*au.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("server panic")
	}
	return &au.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func (s *server) PlaceOrder(ctx context.Context, in *au.PlaceOrderRequest) (*au.PlaceOrderReply, error) {
	aucode := "xxx"
	fmt.Println("PlaceOrder: ", in.AppId, in.AppSecret, in.OrderMsg)
	ok, err := checkAuth(in.AppId, in.AppSecret, aucode)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.Errorf(403, "auth fail", fmt.Sprintf("appid: %s, appsecret: %s", in.AppId, in.AppSecret))
	}

	return &au.PlaceOrderReply{Ok: true}, nil
}
