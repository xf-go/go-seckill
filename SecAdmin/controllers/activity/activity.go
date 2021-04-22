package activity

import (
	"SecAdmin/model"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type ActivityController struct {
	beego.Controller
}

func (p *ActivityController) ListActivity() {
	p.TplName = "activity/list.html"
	p.Layout = "layout/layout.html"

	var err error
	defer func() {
		if err != nil {
			p.Data["Error"] = err
			p.TplName = "activity/error.html"
		}
	}()

	activityModel := model.NewActivityModel()
	activityList, err := activityModel.GetActivityList()
	if err != nil {
		logs.Warn("get activity list failed. err: %v", err)
		return
	}

	p.Data["activity_list"] = activityList
}

func (p *ActivityController) CreateActivity() {
	activityModel := model.NewActivityModel()
	activityList, err := activityModel.GetActivityList()
	if err != nil {
		logs.Warn("get activity list failed. err: %v", err)
		return
	}

	p.Data["activity_list"] = activityList
	p.TplName = "activity/create.html"
	p.Layout = "layout/layout.html"
}
