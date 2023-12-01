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
	UserId    int64     `json:"userid"`
	AuCodes   []string  `json:"aucode"`
	Allow     bool      `json:"allow"`
	Expires   time.Time `json:"expires"`
	// skip time model: created_at, updated_at, deleted_at
}

var ats = []*AccessToken{
	// ww(id:15739) has generated authToken to access 0148P1016_ww
	{Id: 1, AppId: "asdj", AppSecret: "d54sdfejbd561sa", UserId: 15739, AuCodes: []string{"0148P1016_ww"}, Allow: true, Expires: time.Now().Add(time.Hour * 24 * 7)},
	// xjw(id:15743) has generated authToken to access EAMLS1ZT_00
	{Id: 2, AppId: "kfuks", AppSecret: "4fd1ufklnksbry9", UserId: 15743, AuCodes: []string{"EAMLS1ZT_00"}, Allow: true, Expires: time.Now().Add(time.Hour * 24 * 4)},
	// ww(id:15739) has generated authToken to access all au which ww can access.
	{Id: 3, AppId: "jkwsx", AppSecret: "luwxtuf5twprw5l", UserId: 15739, AuCodes: []string{"*"}, Allow: true, Expires: time.Now().Add(time.Hour * 24 * 4)},
	// ww(id:15739) has generated authToken to access all au which ww can access except 0148P1016_ww.
	{Id: 4, AppId: "ggTks", AppSecret: "psuhl055bwaeTIjk", UserId: 15739, AuCodes: []string{"0148P1016_ww"}, Allow: false, Expires: time.Now().Add(time.Hour * 24 * 4)},
	// ww(id:15739) has generated authToken to access [0148P1016_ww, 88853899_ww]
	{Id: 5, AppId: "xstt", AppSecret: "abeo5tgrt754arh57", UserId: 15739, AuCodes: []string{"0148P1016_ww", "88853899_ww"}, Allow: true, Expires: time.Now().Add(time.Hour * 24 * 7)},
	// wsy(id:15747) has generated authToken to access DRW001ZTX_04
	{Id: 5, AppId: "ko8w", AppSecret: "8hw416ery9ah4foig", UserId: 15747, AuCodes: []string{"DRW001ZTX_04"}, Allow: true, Expires: time.Now().Add(time.Hour * 24 * 7)},
	// xjw(id:15743) has generated authToken to access EAMLS1ZT_00
	{Id: 5, AppId: "eut2", AppSecret: "tyt1ra48is13awer6", UserId: 15743, AuCodes: []string{"EAMLS1ZT_00"}, Allow: true, Expires: time.Now().Add(time.Hour * 24 * 7)},
}

func getATByToken(ctx context.Context, appid string, appsecret string, aucode string) (*AccessToken, error) {
	// AppSecret ramains globally unique
	var aimAt *AccessToken
	for _, at := range ats {
		if at.AppId == appid && at.AppSecret == appsecret {
			aimAt = at
			break
		}
	}

	if aimAt == nil {
		return nil, errors.Errorf(403, "Access Token Not Exist", fmt.Sprintf("appid: %s, appsecret: %s", appid, appsecret))
	}

	// only one special case
	if len(aimAt.AuCodes) == 1 && aimAt.AuCodes[0] == "*" {
		return aimAt, nil
	}

	hasAucode := false
	for _, cau := range aimAt.AuCodes {
		if cau == aucode {
			// return aimAt, nil
			hasAucode = true
		}
	}

	// White list or Black list
	if (hasAucode && aimAt.Allow) || (!hasAucode && !aimAt.Allow) {
		return aimAt, nil
	}

	return nil, errors.Errorf(403, "Aucode Is Not Auth", fmt.Sprintf("aucode: %s", aucode))
}
