package controllers

import (
	"SecProxy/service"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type KillController struct {
	web.Controller
}

func (this *KillController) SecKill() {
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
		return
	}

	source := this.GetString("source")
	authcode := this.GetString("authcode")
	secTime := this.GetString("time")
	nance := this.GetString("nance")

	secRequest := &service.SecRequest{}
	secRequest.ProductId = productId
	secRequest.Source = source
	secRequest.AuthCode = authcode
	secRequest.SecTime = secTime
	secRequest.Nance = nance
	secRequest.AccessTime = time.Now()
	secRequest.UserAuthSign = this.Ctx.GetCookie("userAuthSign")
	secRequest.UserId, err = strconv.Atoi(this.Ctx.GetCookie("userId"))
	if err != nil {
		result["code"] = service.ErrInvalidRequest
		result["message"] = fmt.Errorf("invalid cookie:userId").Error()
		return
	}
	if len(this.Ctx.Request.RemoteAddr) > 0 {
		secRequest.ClientAddr = strings.Split(this.Ctx.Request.RemoteAddr, ":")[0]
	}
	secRequest.ClientReferer = this.Ctx.Request.Referer()

	data, code, err := service.SecKill(secRequest)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		return
	}

	result["code"] = code
	result["data"] = data
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
		data, code, err := service.SecInfoList()
		if err != nil {
			result["code"] = 1001
			result["message"] = "invalid product_id"

			logs.Error("invalid request, get product_id failed, err:%v", err)
			return
		}
		result["code"] = code
		result["data"] = data
	} else {

		data, code, err := service.SecInfo(productId)
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()

			logs.Error("invalid request, get product_id failed, err:%v", err)
			return
		}
		result["code"] = code
		result["data"] = data
	}

}
