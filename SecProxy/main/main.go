package main

import (
	_ "SecProxy/router"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

func initConfig() (err error) {
	fmt.Println("web.AppConfig: ", web.AppConfig)
	redisAddr, err := web.AppConfig.String("redis_addr")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	etcdAddr, _ := web.AppConfig.String("etcd_addr")

	secKillConf.RedisAddr = redisAddr
	secKillConf.EtcdAddr = etcdAddr

	logs.Debug("read config succ, redis addr:%v", redisAddr)
	logs.Debug("read config succ, etcd addr:%v", etcdAddr)

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config failed, redis or etcd config is null")
		return
	}

	return
}

func main() {
	err := initConfig()
	if err != nil {
		panic(err)
	}

	err = initSec()
	if err != nil {
		panic(err)
	}

	web.SetStaticPath("/down", "download1")
	web.Run()
}
