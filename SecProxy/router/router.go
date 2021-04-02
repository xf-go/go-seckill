package router

import (
	"SecProxy/controller"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	web.Router("/secinfo", &controller.SkillController{}, "*:SecInfo")
}
