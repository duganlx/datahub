package main

import (
	"net/http"
	"time"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	"gopkg.in/oauth2.v3/errors"
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
	GrantType    string
	ClientId     string
	ClientSecret string
	UserId       int64
}

type AuthManager struct {
	server *oauthServer.Server
}

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
		return "", "", errors.ErrInvalidClient
	}

	return clientId, clientSecret, nil
}
