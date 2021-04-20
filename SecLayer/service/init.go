package service

import "github.com/beego/beego/v2/core/logs"

func InitSecKill(conf *SecLayerConf) (err error) {
	err = initRedis(conf)
	if err != nil {
		logs.Error("init redis failed. err: %v", err)
		return
	}
	logs.Debug("init redis succ.")

	err = initEtcd(conf)
	if err != nil {
		logs.Error("init etcd failed. err: %v", err)
		return
	}
	logs.Debug("init etcd succ.")

	err = loadProductFromEtcd(conf)
	if err != nil {
		logs.Error("load product from etcd failed. err: %v")
		return
	}
	logs.Debug("load product from etcd succ.")

	secLayerContext.secLayerConf = conf
	secLayerContext.Read2HandleChan = make(chan *SecRequest, secLayerContext.secLayerConf.Read2HandleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse, secLayerContext.secLayerConf.Handle2WriteChanSize)
	secLayerContext.HistoryMap = make(map[int]*UserBuyHistory, 100000)
	secLayerContext.productCountMgr = NewProductCountMgr()

	return
}
