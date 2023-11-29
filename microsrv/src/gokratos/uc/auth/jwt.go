package auth

import (
	"encoding/base64"
	"errors"
	"strings"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gopkg.in/oauth2.v3"
)

// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	jwtv4.RegisteredClaims
	Scope string `json:"scope,omitempty"`
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKey    []byte
	SignedMethod jwtv4.SigningMethod
}

func (a *JWTAccessGenerate) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	claims := &JWTAccessClaims{
		RegisteredClaims: jwtv4.RegisteredClaims{
			Issuer:    "",
			Subject:   data.UserID,
			Audience:  jwtv4.ClaimStrings{data.Client.GetID()},
			ExpiresAt: &jwtv4.NumericDate{Time: data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn())},
			NotBefore: &jwtv4.NumericDate{},
			IssuedAt:  &jwtv4.NumericDate{},
			ID:        "",
		},
		Scope: data.TokenInfo.GetScope(),
	}

	token := jwtv4.NewWithClaims(a.SignedMethod, claims)
	var key interface{}
	if a.isECDSA() {
		v, err := jwtv4.ParseECPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isRSAOrPKCS() {
		v, err := jwtv4.ParseRSAPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isHs() {
		key = a.SignedKey
	} else {
		return "", "", errors.New("unsupported sign method")
	}

	access, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}
	return access, refresh, nil
}

func (a *JWTAccessGenerate) isECDSA() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWTAccessGenerate) isRSAOrPKCS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}
