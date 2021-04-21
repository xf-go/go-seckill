package setup

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type MysqlConfig struct {
	UserName string
	Passwd   string
	Port     int
	Database string
	Host     string
}

var AppConf Config

type Config struct {
	MysqlConf MysqlConfig
}

func initConfig() (err error) {
	username, err := beego.AppConfig.String("mysql_user_name")
	if err != nil {
		logs.Error("init config mysql_user_name failed. err: %v", err)
		return
	}
	AppConf.MysqlConf.UserName = username

	pwd, err := beego.AppConfig.String("mysql_passwd")
	if err != nil {
		logs.Error("init config mysql_passwd failed. err: %v", err)
		return
	}
	AppConf.MysqlConf.Passwd = pwd

	host, err := beego.AppConfig.String("mysql_host")
	if err != nil {
		logs.Error("init config mysql_host failed. err: %v", err)
		return
	}
	AppConf.MysqlConf.Host = host

	database, err := beego.AppConfig.String("mysql_database")
	if err != nil {
		logs.Error("init config mysql_database failed. err: %v", err)
		return
	}
	AppConf.MysqlConf.Database = database

	port, err := beego.AppConfig.Int("mysql_port")
	if err != nil {
		logs.Error("init config mysql_port failed. err: %v", err)
		return
	}
	AppConf.MysqlConf.Port = port

	return
}
