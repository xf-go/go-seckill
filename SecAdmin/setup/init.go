package setup

import (
	"fmt"

	"SecAdmin/model"

	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Init() (err error) {
	err = initConfig()
	if err != nil {
		logs.Warn("init config failed. err: %v", err)
		return
	}

	db, err := initDB()
	if err != nil {
		logs.Warn("init db failed. err: %v", err)
		return
	}

	model.Init(db)

	return
}

func initDB() (db *sqlx.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		AppConf.MysqlConf.UserName,
		AppConf.MysqlConf.Passwd,
		AppConf.MysqlConf.Host,
		AppConf.MysqlConf.Port,
		AppConf.MysqlConf.Database,
	)
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		logs.Error("open mysql failed. err: %v", err)
		return
	}
	logs.Debug("connect to mysql succ.")
	return
}
