package router

import (
	"SecProxy/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/seckill", &controllers.KillController{}, "*:SecKill")
	beego.Router("/secinfo", &controllers.KillController{}, "*:SecInfo")
}
