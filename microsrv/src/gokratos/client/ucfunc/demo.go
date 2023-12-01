package ucfunc

import "fmt"

// Demo1 用户ww生成*只能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`0148P1016_ww`
func Demo1(connType string) {
	token, err := loginTemplate("asdj", "d54sdfejbd561sa", "0148P1016_ww", connType)
	if err != nil {
		fmt.Printf("[Demo1 %s] X: %v\n", connType, err)
		return
	}

	fmt.Printf("[Demo1 %s] √: %v\n", connType, token)
}

// Demo2 用户ww生成*只能*访问资产单元`[0148P1016_ww, 88853899_ww]`的访问令牌，并访问资产单元`0148P1016_ww`
func Demo2(connType string) {
	token, err := loginTemplate("xstt", "abeo5tgrt754arh57", "0148P1016_ww", connType)
	if err != nil {
		fmt.Printf("[Demo2 %s] X: %v\n", connType, err)
		return
	}

	fmt.Printf("[Demo2 %s] √: %v\n", connType, token)
}

// Demo3 用户ww生成*只能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`88853899_ww`
func Demo3(connType string) {
	token, err := loginTemplate("asdj", "d54sdfejbd561sa", "88853899_ww", connType)
	if err != nil {
		fmt.Printf("[Demo3 %s] √: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo3 %s] X: %v\n", connType, token)
}

// Demo4 用户ww生成可访问*所有*资产单元（`MANAGER_WW组`）的访问令牌，并访问资产单元`88853899_ww`
func Demo4(connType string) {
	token, err := loginTemplate("jkwsx", "luwxtuf5twprw5l", "88853899_ww", connType)
	if err != nil {
		fmt.Printf("[Demo4 %s] X: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo4 %s] √: %v\n", connType, token)
}

// Demo5 用户ww生成可访问*所有*资产单元（`MANAGER_WW组`）的访问令牌，并访问资产单元`EAMLS1ZT_00`
func Demo5(connType string) {
	token, err := loginTemplate("jkwsx", "luwxtuf5twprw5l", "EAMLS1ZT_00", connType)
	if err != nil {
		fmt.Printf("[Demo5 %s] √: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo5 %s] X: %v\n", connType, token)
}

// Demo6 用户ww生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`0148P1016_ww`
func Demo6(connType string) {
	token, err := loginTemplate("ggTks", "psuhl055bwaeTIjk", "0148P1016_ww", connType)
	if err != nil {
		fmt.Printf("[Demo6 %s] √: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo6 %s] X: %v\n", connType, token)
}

// Demo7 用户ww生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`88853899_ww`
func Demo7(connType string) {
	token, err := loginTemplate("ggTks", "psuhl055bwaeTIjk", "88853899_ww", connType)
	if err != nil {
		fmt.Printf("[Demo7 %s] X: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo7 %s] √: %v\n", connType, token)
}

// Demo8 用户ww生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`EAMLS1ZT_00`
func Demo8(connType string) {
	token, err := loginTemplate("ggTks", "psuhl055bwaeTIjk", "EAMLS1ZT_00", connType)
	if err != nil {
		fmt.Printf("[Demo8 %s] √: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo8 %s] X: %v\n", connType, token)
}

// Demo9 用户wsy生成*只能*访问资产单元`[DRW001ZTX_04]`的访问令牌，并访问资产单元`DRW001ZTX_04`
func Demo9(connType string) {
	token, err := loginTemplate("ko8w", "8hw416ery9ah4foig", "DRW001ZTX_04", connType)
	if err != nil {
		fmt.Printf("[Demo9 %s] X: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo9 %s] √: %v\n", connType, token)
}

// Demo10 用户xjw生成*只能*访问资产单元`[EAMLS1ZT_00]`的访问令牌，并访问资产单元`EAMLS1ZT_00`
func Demo10(connType string) {
	token, err := loginTemplate("eut2", "tyt1ra48is13awer6", "EAMLS1ZT_00", connType)
	if err != nil {
		fmt.Printf("[Demo10 %s] X: %v\n", connType, err.Error())
		return
	}

	fmt.Printf("[Demo10 %s] √: %v\n", connType, token)
}
