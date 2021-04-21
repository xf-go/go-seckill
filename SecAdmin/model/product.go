package model

import (
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

type ProductModel struct {
}

type Product struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Total  int    `db:"total"`
	Status int    `db:"status"`
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) GetProductList() (list []*Product, err error) {
	sql := "select id, name, total, status from product"
	err = DB.Select(&list, sql)
	if err != nil {
		logs.Error("select from mysql failed, err: %v, sql: %v", err, sql)
		return
	}
	return
}

func (p *ProductModel) CreateProduct(product *Product) (err error) {
	sql := "insert into product(name, total, status) values (?,?,?)"
	_, err = DB.Exec(sql, product.Name, product.Total, product.Status)
	if err != nil {
		logs.Error("create product failed, err: %v, sql: %v", err, sql)
		return
	}
	logs.Debug("insert into database succ.")
	return
}
