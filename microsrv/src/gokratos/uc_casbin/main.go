package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

func main() {
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println(e)
}
