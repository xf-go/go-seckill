package product

import (
	"SecAdmin/model"
	_ "SecAdmin/model"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {
	p.TplName = "product/list.html"
	p.Layout = "layout/layout.html"

	var err error
	defer func() {
		if err != nil {
			p.Data["Error"] = err
			p.TplName = "product/error.html"
		}
	}()

	productModel := model.NewProductModel()
	productList, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed. err: %v", err)
		return
	}

	p.Data["product_list"] = productList

}

func (p *ProductController) CreateProduct() {
	productModel := model.NewProductModel()
	productList, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed. err: %v", err)
		return
	}

	p.Data["product_list"] = productList
	p.TplName = "product/create.html"
	p.Layout = "layout/layout.html"
}

func (p *ProductController) SubmitProduct() {
	p.TplName = "product/create.html"
	p.Layout = "layout/layout.html"
	errNsg := "success"

	var err error
	defer func() {
		if err != nil {
			p.Data["Error"] = errNsg
			p.TplName = "product/error.html"
		}
	}()

	productName := p.GetString("product_name")
	productTotal, err := p.GetInt("product_total")
	if err != nil {
		logs.Warn("invalid product total, err: %v", err)
		errNsg = fmt.Sprintf("invalid product total, err: %v", err)
		return
	}
	productStatus, err := p.GetInt("product_status")
	if err != nil {
		logs.Warn("invalid product status, err: %v", err)
		errNsg = fmt.Sprintf("invalid product status, err: %v", err)
		return
	}

	productModel := model.NewProductModel()
	product := model.Product{
		Name:   productName,
		Total:  productTotal,
		Status: productStatus,
	}
	err = productModel.CreateProduct(&product)
	if err != nil {
		logs.Warn("create product failed. err: %v", err)
		errNsg = fmt.Sprintf("create product failed. err: %v", err)
		return
	}

	logs.Debug("product name[%s], product total[%d], product status[%d]", productName, productTotal, productStatus)
}
