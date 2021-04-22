package router

import (
	"SecAdmin/controllers/activity"
	"SecAdmin/controllers/product"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/product/list", &product.ProductController{}, "*:ListProduct")
	beego.Router("/product/create", &product.ProductController{}, "*:CreateProduct")
	beego.Router("/product/submit", &product.ProductController{}, "*:SubmitProduct")

	beego.Router("/activity/list", &activity.ActivityController{}, "*:ListActivity")
	beego.Router("/activity/create", &activity.ActivityController{}, "*:CreateActivity")
}
