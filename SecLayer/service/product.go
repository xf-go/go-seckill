package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/core/logs"
	etcd "go.etcd.io/etcd/client/v3"
)

func loadProductFromEtcd(conf *SecLayerConf) (err error) {
	logs.Debug("start load product from etcd.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := secLayerContext.etcdClient.Get(ctx, conf.EtcdConfig.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err: %v", conf.EtcdConfig.EtcdSecProductKey, err)
		return
	}
	logs.Debug("load product from etcd srcc. resp: %v", resp)

	var secProductInfo []SecProductInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] value[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("json.Unmarshal failed, err: %v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	updateSecProductInfo(conf, secProductInfo)
	initSecProductWatcher(conf)

	return
}

func updateSecProductInfo(conf *SecLayerConf, secProductInfo []SecProductInfoConf) {
	var tmp map[int]*SecProductInfoConf = make(map[int]*SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		productInfo := v
		productInfo.secLimit = &SecLimit{}
		tmp[v.ProductId] = &productInfo
	}
	secLayerContext.RWSecProductLock.Lock()
	conf.SecProductInfoMap = tmp
	secLayerContext.RWSecProductLock.Unlock()
}

func initSecProductWatcher(conf *SecLayerConf) {
	go watchSecProductKey(conf)
}

func watchSecProductKey(conf *SecLayerConf) {
	key := conf.EtcdConfig.EtcdSecProductKey
	logs.Debug("begin watch key:%s", key)
	var err error
	for {
		ch := secLayerContext.etcdClient.Watch(context.Background(), key)
		var secProductInfo []SecProductInfoConf
		getConfSucc := true

		for v := range ch {
			for _, ev := range v.Events {
				if ev.Type == etcd.EventTypeDelete {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if ev.Type == etcd.EventTypePut && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], json.Unmarshal failed, err:%v", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd,%s,$q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(conf, secProductInfo)
			}
		}
	}
}
