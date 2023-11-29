package server

import (
	"context"
	"fmt"
	uc "gokratos/api/uc/v1"
	"gokratos/uc/auth"

	errors "github.com/go-kratos/kratos/v2/errors"
)

func (s *Server) Login(ctx context.Context, in *uc.LoginRequest) (*uc.LoginReply, error) {
	if in.AccessType != "code" {
		return nil, errors.New(303, "unsupport access type", fmt.Sprintf("unsupport type %s", in.AccessType))
	}
	if len(in.AppId) == 0 || len(in.AppSecret) == 0 {
		return nil, errors.Errorf(400, "appid and appsecret should not be null", fmt.Sprintf("appid: %s, appsecret: %s", in.AppId, in.AppSecret))
	}

	user, err := getUserByCode(ctx, in.AppId, in.AppSecret)
	if err != nil {
		return nil, err
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
