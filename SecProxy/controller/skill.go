package controller

import (
	"github.com/beego/beego/v2/server/web"
)

type SkillController struct {
	web.Controller
}

func (this *SkillController) SecKill() {
	this.Data["json"] = "sec kill"
	this.ServeJSON()
}

func (this *SkillController) SecInfo() {
	this.Data["json"] = "sec info"
	this.ServeJSON()
}
