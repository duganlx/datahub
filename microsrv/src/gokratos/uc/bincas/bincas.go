package bincas

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

type CBEnforce struct {
	cbe *casbin.Enforcer
}

func NewCasbinEnforcer() (*CBEnforce, error) {
	e, err := casbin.NewEnforcer("./bincas/model.conf", "./bincas/policy.csv")

	return &CBEnforce{
		cbe: e,
	}, err
}

func (e CBEnforce) CanAccessAu(sub, aucode, opType string) (bool, error) {
	obj := fmt.Sprintf("AU_%s", aucode)
	ok, err := e.cbe.Enforce(sub, obj, opType)
	return ok, err
}

func (e CBEnforce) UpdateAuth() {
	fmt.Println("权限更新...")
	// 给部门 it 访问 121000 操作类型 w 的权限
	e.cbe.AddPolicy("DEPT_it", "AU_121000", "w")
	// 删除 用户 wsy 访问 DRW001ZTX_04 的操作类型 * 的权限
	e.cbe.RemovePolicy("USER_wsy", "AU_DRW001ZTX_04", "*")
	// 删除 用户 yrl 投资经理 MANAGER_WW 的权限
	e.cbe.DeleteRoleForUser("USER_yrl", "MANAGER_WW")

	// todo
}
