package service

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	etcd "go.etcd.io/etcd/client/v3"
)

func initEtcd(conf *SecLayerConf) (err error) {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{conf.EtcdConfig.EtcdAddr},
		DialTimeout: time.Duration(conf.EtcdConfig.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err: %v", err)
		return
	}

	secLayerContext.etcdClient = cli
	return
}
