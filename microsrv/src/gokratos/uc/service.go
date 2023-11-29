package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	uc "gokratos/api/uc/v1"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"strings"

// 	errors "github.com/go-kratos/kratos/v2/errors"
// )

// type server struct {
// 	uc.UnimplementedUserCenterServer
// }

// func (s *server) Login(ctx context.Context, in *uc.LoginRequest) (*uc.LoginReply, error) {
// 	// 参数校验
// 	// == begin ==
// 	if in.AccessType != "code" {
// 		return nil, errors.New(303, "unsupport access type", fmt.Sprintf("unsupport type %s", in.AccessType))
// 	}
// 	if len(in.AppId) == 0 || len(in.AppSecret) == 0 {
// 		return nil, errors.Errorf(400, "appid and appsecret should not be null", fmt.Sprintf("appid: %s, appsecret: %s", in.AppId, in.AppSecret))
// 	}
// 	// == end ==

// 	user, err := s.getUserByCode(ctx, in.AppId, in.AppSecret)
// 	if err != nil {
// 		return nil, err
// 	}

// 	auth, err := s.getToken(ctx, &TokenRequest{
// 		GrantType:    "password",
// 		ClientId:     ClientId,
// 		ClientSecret: ClientSecret,
// 		UserId:       user.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &uc.LoginReply{
// 		AccessToken:  auth.AccessToken,
// 		RefreshToken: auth.RefreshToken,
// 		TokenType:    auth.TokenType,
// 		Expires:      auth.ExpiresIn,
// 		Scrope:       "",
// 		Uid:          user.Id,
// 	}, nil
// }

// func (s *server) getUserByCode(ctx context.Context, appid string, appsecret string) (*User, error) {
// 	if appid == AppId && appsecret == AppSecret {
// 		return &User{Id: int64(UserId), UserName: UserName, Mobile: UserMobile}, nil
// 	}

// 	return nil, errors.Errorf(403, "Authentication failed", fmt.Sprintf("appid: %s, appsecret: %s", appid, appsecret))
// }

// func (s *server) getToken(ctx context.Context, req *TokenRequest) (*Auth, error) {
// 	authMgr := NewAuthManager()

// 	// 构造并发送到 OAuth2 服务器验证的请求
// 	v := make(url.Values)
// 	v.Set("grant_type", req.GrantType)
// 	v.Set("client_id", req.ClientId)
// 	v.Set("client_secret", req.ClientSecret)
// 	v.Set("username", strconv.FormatInt(req.UserId, 10))
// 	v.Set("password", "x")
// 	httpReq, err := http.NewRequest("POST", "/", strings.NewReader(v.Encode()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

// 	gt, tgr, err := authMgr.server.ValidationTokenRequest(httpReq)
// 	if err != nil {
// 		fmt.Printf("ERROR token ValidationTokenRequest: %+v\n", httpReq)
// 		return nil, errors.Errorf(500, "inValid Token Request", err.Error())
// 	}

// 	ti, err := authMgr.server.GetAccessToken(gt, tgr)
// 	if err != nil {
// 		fmt.Printf("ERROR gt: %+v, tgr: %+v", gt, tgr)
// 		return nil, errors.Errorf(500, "GetAccessToken Failed", err.Error())
// 	}

// 	data := authMgr.server.GetTokenData(ti)
// 	// fmt.Printf("DEBUG token data info: %+v\n", data)

// 	buf, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Printf("ERROR Marshal data: %+v", data)
// 		return nil, errors.New(500, "Json Marshal Failed", err.Error())
// 	}
// 	var auth Auth
// 	if err := json.Unmarshal(buf, &auth); err != nil {
// 		fmt.Printf("ERROR Unmarshal buffer: %s", string(buf))
// 		return nil, errors.New(500, "Unmarshal buffer Failed", err.Error())
// 	}

// 	// fmt.Printf("DEBUG generate auth info: %+v\n", auth)
// 	return &auth, nil
// }
