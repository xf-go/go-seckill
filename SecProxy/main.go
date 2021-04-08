package main

import (
	_ "SecProxy/router"
	"SecProxy/setup"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	err := setup.InitConfig()
	if err != nil {
		panic(err)
	}

	err = setup.InitSec()
	if err != nil {
		panic(err)
	}

	beego.SetStaticPath("/down", "download1")
	beego.Run()
}
