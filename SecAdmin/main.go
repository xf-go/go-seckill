package main

import (
	"fmt"

	_ "SecAdmin/router"
	"SecAdmin/setup"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	err := setup.Init()
	if err != nil {
		panic(fmt.Sprintf("init failed, err: %v", err))
	}
	logs.Debug("run.")
	beego.Run()
}
