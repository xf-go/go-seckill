package setup

import (
	"SecProxy/service"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	etcd "go.etcd.io/etcd/clientv3"
)

var (
	redisPool  *redis.Pool
	etcdClient *etcd.Client
)

func InitSec() (err error) {
	err = initLogger()
	if err != nil {
		logs.Error("init logger failed, err: %v", err)
		return
	}

	// err = initRedis()
	// if err != nil {
	// 	logs.Error("init redis failed, err: %v", err)
	// 	return
	// }

	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed, err: %v", err)
		return
	}

	err = loadSecConf()
	if err != nil {
		logs.Error("load sec conf failed, err: %v", err)
		return
	}

	service.InitService(secKillServer)
	initSecProductWatcher()

	logs.Info("init sec succ")
	return
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = secKillServer.LogPath
	config["level"] = convertLogLevel(secKillServer.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err: ", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))

	return
}

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initRedis() (err error) {
	redisPool = &redis.Pool{
		MaxIdle:     secKillServer.RedisBlackConf.RedisMaxIdle,
		MaxActive:   secKillServer.RedisBlackConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillServer.RedisBlackConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillServer.RedisBlackConf.RedisAddr)
		},
	}
	conn := redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err: %v", err)
		return
	}

	return
}

func initEtcd() (err error) {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{secKillServer.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(secKillServer.EtcdConf.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err: %v", err)
		return
	}
	etcdClient = cli
	return
}

func loadSecConf() (err error) {
	resp, err := etcdClient.Get(context.Background(), secKillServer.EtcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err: %v", secKillServer.EtcdConf.EtcdSecProductKey, err)
		return
	}

	var secProductInfo []service.SecProductInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] value[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("json.Unmarshal failed, err: %v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	updateSecProductInfo(secProductInfo)
	return
}

func initSecProductWatcher() {
	go watchSecProductKey(secKillServer.EtcdConf.EtcdSecProductKey)
}

func watchSecProductKey(key string) {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:%v", err)
		return
	}
	logs.Debug("begin watch key:%s", key)
	for {
		ch := cli.Watch(context.Background(), key)
		var secProductInfo []service.SecProductInfoConf
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
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}

func updateSecProductInfo(secProductInfo []service.SecProductInfoConf) {
	var tmp map[int]*service.SecProductInfoConf = make(map[int]*service.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		productInfo := v
		tmp[v.ProductId] = &productInfo
	}
	secKillServer.RWSecProductLock.Lock()
	secKillServer.SecProductInfoMap = tmp
	secKillServer.RWSecProductLock.Unlock()
}
