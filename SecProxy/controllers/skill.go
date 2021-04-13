package controllers

import (
	"SecProxy/service"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type KillController struct {
	web.Controller
}

func (this *KillController) SecKill() {
	this.Data["json"] = "sec kill"
	this.ServeJSON()
}

func (this *KillController) SecInfo() {
	productId, err := this.GetInt("product_id")
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "success"
	defer func() {
		this.Data["json"] = result
		this.ServeJSON()
	}()
	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid product_id"

		logs.Error("invalid request, get product_id failed, err:%v", err)
		return
	}

	data, code, err := service.SecInfo(productId)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()

		logs.Error("invalid request, get product_id failed, err:%v", err)
		return
	}
	result["data"] = data
}
