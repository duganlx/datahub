package bincas

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

func CanAccessAu(uid, aucode string) (bool, error) {
	e, err := casbin.NewEnforcer("./bincas/model.conf", "./bincas/policy.csv")
	if err != nil {
		return false, err
	}

	sub := fmt.Sprintf("USER_%s", uid)
	obj := fmt.Sprintf("AU_%s", aucode)

	ok, err := e.Enforce(sub, obj)
	return ok, err
}
