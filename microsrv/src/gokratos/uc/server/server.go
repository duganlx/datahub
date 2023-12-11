package server

import (
	"context"
	"fmt"
	uc "gokratos/api/uc/v1"
	"gokratos/uc/auth"
	"gokratos/uc/bincas"

	errors "github.com/go-kratos/kratos/v2/errors"
)

type Server struct {
	uc.UnimplementedUserCenterServer

	Cbe *bincas.CBEnforce
}

func NewServer() (*Server, error) {
	cbe, err := bincas.NewCasbinEnforcer()
	if err != nil {
		return nil, err
	}

	return &Server{
		Cbe: cbe,
	}, nil
}

func (s *Server) Login(ctx context.Context, in *uc.LoginRequest) (*uc.LoginReply, error) {
	if in.AccessType != "code" {
		return nil, errors.New(303, "unsupport access type", fmt.Sprintf("unsupport type %s", in.AccessType))
	}
	if len(in.AppId) == 0 || len(in.AppSecret) == 0 {
		return nil, errors.Errorf(400, "appid and appsecret should not be null", fmt.Sprintf("appid: %s, appsecret: %s", in.AppId, in.AppSecret))
	}

	user, err := getUserByCode(ctx, in.AppId, in.AppSecret, in.AuCode)
	if err != nil {
		return nil, err
	}

	ok, err := s.Cbe.CanAccessAu(user.UserName, in.AuCode, "*")
	if err != nil {
		return nil, errors.Errorf(500, "Server Internal Error", err.Error())
	}
	if !ok {
		return nil, errors.Errorf(403, "Access AuCode Failed", fmt.Sprintf("user: %+v, auCode: %s", user, in.AuCode))
	}

	authMgr := auth.NewAuthManager()
	auth, err := authMgr.GetToken(ctx, &auth.TokenRequest{
		GrantType: "password",
		UserId:    user.Id,
	})
	if err != nil {
		return nil, err
	}

	return &uc.LoginReply{
		AccessToken:  auth.AccessToken,
		RefreshToken: auth.RefreshToken,
		TokenType:    auth.TokenType,
		Expires:      auth.ExpiresIn,
		Scrope:       auth.Scope,
		Uid:          user.Id,
	}, nil
}

func (s *Server) Authraw(ctx context.Context, in *uc.AuthrawRequest) (*uc.AuthrawReply, error) {
	user, usertype, err := getUserById(ctx, in.Uid)
	if err != nil {
		return nil, errors.Errorf(500, "Server Internal Error", err.Error())
	}

	var sub string = ""
	switch usertype {
	case "user":
		sub = fmt.Sprintf("USER_%s", user.UserName)
	case "dept":
		sub = fmt.Sprintf("DEPT_%s", user.UserName)
	default:
	}

	ok, err := s.Cbe.CanAccessAu(sub, in.AuCode, in.OpType)
	if err != nil {
		return nil, errors.Errorf(500, "Server Internal Error", err.Error())
	}

	return &uc.AuthrawReply{
		Ok: ok,
	}, nil
}
