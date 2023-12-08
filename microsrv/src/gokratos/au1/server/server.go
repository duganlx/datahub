package server

import (
	"context"
	"fmt"
	au "gokratos/api/au/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

type server struct {
	au.UnimplementedAssetUnitServer

	aucode string
	akpair map[string]string // appsecret: appid
}

func NewAuServer(aucode string) *server {
	s := &server{
		aucode: aucode,
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
	var ok bool = false
	var err error

	// check cache
	if s.hasAkpair(in.AppId, in.AppSecret) {
		ok = true
	} else {
		ok, err = checkAuth(ctx, in.AppId, in.AppSecret, s.aucode)

		if err != nil {
			return nil, err
		}
	}

	if !ok {
		return nil, errors.Errorf(403, "auth fail", fmt.Sprintf("appid: %s, appsecret: %s", in.AppId, in.AppSecret))
	}
	
	s.addAkpair(in.AppId, in.AppSecret)

	fmt.Println("PlaceOrder: ", in.OrderMsg)

	return &au.PlaceOrderReply{Ok: true}, nil
}

func (s *server) hasAkpair(appid, appsecret string) bool {
	as, exist := s.akpair[appid]

	if !exist {
		return false
	}

	return as == appsecret
}

func (s *server) addAkpair(appid, appsecret string) {
	exist := s.hasAkpair(appid, appsecret)

	if exist {
		return
	}

	s.akpair[appid] = appsecret
}