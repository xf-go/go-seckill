package main

import (
	"fmt"

	"SecLayer/service"
	"SecLayer/setup"

	"github.com/beego/beego/v2/core/logs"
)

func main() {
	// 加载配置文件
	err := setup.InitConfig("ini", "./conf/secLayer.conf")
	if err != nil {
		logs.Error("init config failed. err: %v", err)
		panic(err)
	}

	// 初始化日志库
	err = setup.InitLogger()
	if err != nil {
		logs.Error("init logger failed. err:%v", err)
		panic(err)
	}

	// 初始化秒杀逻辑
	err = service.InitSecKill(setup.AppConfig)
	if err != nil {
		msg := fmt.Sprintf("init sec kill failed. err: %v", err)
		logs.Error(msg)
		panic(msg)
	}

	// 运行业务逻辑
	err = service.Run()
	if err != nil {
		logs.Error("service run failed. err: %v", err)
		return
	}

	logs.Info("service run exited.")
}
