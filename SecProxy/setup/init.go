package setup

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
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

	initSecProductWatcher()

	logs.Info("init sec succ")
	return
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = secKillConf.LogPath
	config["level"] = convertLogLevel(secKillConf.LogLevel)

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
		MaxIdle:     secKillConf.RedisAddr.RedisMaxIdle,
		MaxActive:   secKillConf.RedisAddr.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisAddr.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisAddr.RedisAddr)
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
		Endpoints:   []string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err: %v", err)
		return
	}
	etcdClient = cli
	return
}

func loadSecConf() (err error) {
	resp, err := etcdClient.Get(context.Background(), secKillConf.EtcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err: %v", secKillConf.EtcdConf.EtcdSecProductKey, err)
		return
	}

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

	updateSecProductInfo(secProductInfo)
	return
}

func initSecProductWatcher() {
	go watchSecProductKey(secKillConf.EtcdConf.EtcdSecProductKey)
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
		var secProductInfo []SecProductInfoConf
		getConfSucc := true

		for v := range ch {
			for _, event := range v.Events {
				if event.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if event.Type == mvccpb.PUT && string(event.Kv.Key) == key {
					err = json.Unmarshal(event.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], json.Unmarshal failed, err:%v", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd,%s,$q : %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}

func updateSecProductInfo(secProductInfo []SecProductInfoConf) {
	var tmp map[int]*SecProductInfoConf = make(map[int]*SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		tmp[v.ProductID] = &v
	}
	secKillConf.rwSecProductLock.Lock()
	secKillConf.SecProductInfoMap = tmp
	secKillConf.rwSecProductLock.Unlock()
}
