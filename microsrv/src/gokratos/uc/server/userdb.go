package server

import (
	"context"
	"fmt"

	errors "github.com/go-kratos/kratos/v2/errors"
)

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	Mobile   string `json:"mobile"`
}

var users = []*User{
	{Id: 15739, UserName: "ww", Mobile: "15308681364"},
	{Id: 15743, UserName: "xjw", Mobile: "13608681364"},
	{Id: 15747, UserName: "wsy", Mobile: "13708681364"},
}

func getUserByCode(ctx context.Context, appid string, appsecret string, aucode string) (*User, error) {
	at, err := getATByToken(ctx, appid, appsecret, aucode)
	if err != nil {
		return nil, err
	}

	uid := at.UserId
	for _, u := range users {
		if u.Id == uid {
			return u, nil
		}
	}

	return nil, errors.Errorf(404, "User Not Found", fmt.Sprintf("uid: %d", uid))
}
