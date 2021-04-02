package main

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
)

func initSec() (err error) {
	err = initRedis()
	if err != nil {
		logs.Error("init redis failed, err:%v", err)
		return
	}

	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}

	return
}

func initRedis() (err error) {
	pool := redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisAddr)
		},
	}
	fmt.Println("pool: ", pool)
	return
}

func initEtcd() (err error) {
	return
}
