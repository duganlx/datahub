package server

import (
	"context"
	"fmt"
	"time"

	errors "github.com/go-kratos/kratos/v2/errors"
)

type AccessToken struct {
	Id        int64     `json:"id"`
	AppId     string    `json:"appid"`
	AppSecret string    `json:"appsecret"`
	Expires   time.Time `json:"expires"`
	UserId    int64     `json:"userid"`
}

var ats = []*AccessToken{
	{Id: 1, AppId: "asdj", AppSecret: "d54sdfejbd561sa", Expires: time.Now().Add(time.Hour * 24 * 7), UserId: 15739},
	{Id: 2, AppId: "kfuks", AppSecret: "4fd1ufklnksbry9", Expires: time.Now().Add(time.Hour * 24 * 4), UserId: 15743},
}

func getATByToken(ctx context.Context, appid string, appsecret string) (*AccessToken, error) {

	for _, at := range ats {
		if at.AppId == appid && at.AppSecret == appsecret {
			return at, nil
		}
	}

	return nil, errors.Errorf(403, "Authentication failed", fmt.Sprintf("appid: %s, appsecret: %s", appid, appsecret))
}
