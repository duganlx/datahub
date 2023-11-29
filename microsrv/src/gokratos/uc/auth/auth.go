package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	errors "github.com/go-kratos/kratos/v2/errors"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	oerrors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	oauthServer "gopkg.in/oauth2.v3/server"
)

type Auth struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type TokenRequest struct {
	GrantType string
	UserId    int64
}

type AuthManager struct {
	server *oauthServer.Server
}

var (
	ClientId     = "thisisaclientid"
	ClientSecret = "thisisaclientsecret"
)

func NewAuthManager() *AuthManager {
	authMgr := &AuthManager{}
	manager := manage.NewDefaultManager()
	clientStore := &ClientStore{}
	tokenStore := &TokenStore{}
	accessGenerate := &JWTAccessGenerate{
		SignedKey:    []byte("thisisaapikey"),
		SignedMethod: jwtv4.SigningMethodHS256,
	}

	manager.SetPasswordTokenCfg(
		&manage.Config{
			AccessTokenExp:    time.Hour * 24 * 7,
			RefreshTokenExp:   time.Hour * 24 * 10,
			IsGenerateRefresh: true,
		},
	)
	manager.MapClientStorage(clientStore)
	manager.MapTokenStorage(tokenStore)
	manager.MapAccessGenerate(accessGenerate)

	srv := oauthServer.NewServer(oauthServer.NewConfig(), manager)
	srv.SetClientInfoHandler(authMgr.ClientFormHandler)
	srv.SetPasswordAuthorizationHandler(authMgr.PasswordAuthorizationHandler)

	authMgr.server = srv

	return authMgr
}

func (am *AuthManager) PasswordAuthorizationHandler(account, pwd string) (string, error) {
	return account, nil
}

func (am *AuthManager) ClientFormHandler(r *http.Request) (string, string, error) {
	clientId := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	if clientId == "" || clientSecret == "" {
		return "", "", oerrors.ErrInvalidClient
	}

	return clientId, clientSecret, nil
}

func (am *AuthManager) GetToken(ctx context.Context, req *TokenRequest) (*Auth, error) {
	// 构造并发送到 OAuth2 服务器验证的请求
	v := make(url.Values)
	v.Set("grant_type", req.GrantType)
	v.Set("client_id", ClientId)
	v.Set("client_secret", ClientSecret)
	v.Set("username", strconv.FormatInt(req.UserId, 10))
	v.Set("password", "x")
	httpReq, err := http.NewRequest("POST", "/", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	gt, tgr, err := am.server.ValidationTokenRequest(httpReq)
	if err != nil {
		fmt.Printf("ERROR token ValidationTokenRequest: %+v\n", httpReq)
		return nil, errors.Errorf(500, "inValid Token Request", err.Error())
	}

	ti, err := am.server.GetAccessToken(gt, tgr)
	if err != nil {
		fmt.Printf("ERROR gt: %+v, tgr: %+v", gt, tgr)
		return nil, errors.Errorf(500, "GetAccessToken Failed", err.Error())
	}

	data := am.server.GetTokenData(ti)
	// fmt.Printf("DEBUG token data info: %+v\n", data)

	buf, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("ERROR Marshal data: %+v", data)
		return nil, errors.New(500, "Json Marshal Failed", err.Error())
	}
	var auth Auth
	if err := json.Unmarshal(buf, &auth); err != nil {
		fmt.Printf("ERROR Unmarshal buffer: %s", string(buf))
		return nil, errors.New(500, "Unmarshal buffer Failed", err.Error())
	}

	// fmt.Printf("DEBUG generate auth info: %+v\n", auth)
	return &auth, nil
}
