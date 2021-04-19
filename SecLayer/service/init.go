package service

import "github.com/beego/beego/v2/core/logs"

func InitSecKill(conf *SecLayerConf) (err error) {
	err = initRedis(conf)
	if err != nil {
		logs.Error("init redis failed. err: %v", err)
		return
	}

	err = initEtcd(conf)
	if err != nil {
		logs.Error("init etcd failed. err: %v", err)
		return
	}

	err = loadProductFromEtcd(conf)
	if err != nil {
		logs.Error("load product from etcd failed. err: %v")
		return
	}

	secLayerContext.secLayerConf = conf
	secLayerContext.Read2HandleChan = make(chan *SecRequest, secLayerContext.secLayerConf.Read2HandleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse, secLayerContext.secLayerConf.Handle2WriteChanSize)
	return
}
