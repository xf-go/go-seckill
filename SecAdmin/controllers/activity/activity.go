package activity

import (
	"SecAdmin/model"
	"fmt"
	"net/http"

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

func (p *ActivityController) SubmitActivity() {
	p.TplName = "activity/create.html"
	p.Layout = "layout/layout.html"
	errNsg := "success"

	var err error
	defer func() {
		if err != nil {
			p.Data["Error"] = errNsg
			p.TplName = "activity/error.html"
		}
	}()

	activityName := p.GetString("activity_name")
	if len(activityName) == 0 {
		err = fmt.Errorf("invalid activity total, err: %v", err)
		errNsg = "活动名字不能为空"
		return
	}
	productId, err := p.GetInt("product_id")
	if err != nil {
		err = fmt.Errorf("invalid product id, err: %v", err)
		errNsg = err.Error()
		return
	}
	startTime, err := p.GetInt64("start_time")
	if err != nil {
		err = fmt.Errorf("invalid start time, err: %v", err)
		errNsg = err.Error()
		return
	}
	endTime, err := p.GetInt64("end_time")
	if err != nil {
		err = fmt.Errorf("invalid end time, err: %v", err)
		errNsg = err.Error()
		return
	}
	total, err := p.GetInt64("total")
	if err != nil {
		err = fmt.Errorf("invalid total, err: %v", err)
		errNsg = err.Error()
		return
	}
	speed, err := p.GetInt("speed")
	if err != nil {
		err = fmt.Errorf("invalid speed, err: %v", err)
		errNsg = err.Error()
		return
	}
	limit, err := p.GetInt("buy_limit")
	if err != nil {
		err = fmt.Errorf("invalid buy_limit, err: %v", err)
		errNsg = err.Error()
		return
	}

	activityModel := model.NewActivityModel()
	activity := model.Activity{
		Name:      activityName,
		ProductId: productId,
		StartTime: startTime,
		EndTime:   endTime,
		Total:     total,
		Speed:     speed,
		BuyLimit:  limit,
	}
	err = activityModel.CreateActivity(&activity)
	if err != nil {
		err = fmt.Errorf("create activity failed. err: %v", err)
		errNsg = err.Error()
		return
	}

	p.Redirect("/activity/list", http.StatusFound)
}
