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
	{Id: 15739, UserName: "Tom", Mobile: "15308681364"},
	{Id: 15743, UserName: "Jim", Mobile: "13608681364"},
}

func getUserByCode(ctx context.Context, appid string, appsecret string) (*User, error) {
	at, err := getATByToken(ctx, appid, appsecret)
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
