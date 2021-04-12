package setup

import (
	"fmt"
	"strings"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

var (
	secKillConf = &SecKillConf{
		SecProductInfoMap: make(map[int]*SecProductInfoConf, 1024),
	}
)

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecKillConf struct {
	RedisAddr         RedisConf
	EtcdConf          EtcdConf
	LogPath           string
	LogLevel          string
	SecProductInfoMap map[int]*SecProductInfoConf
	rwSecProductLock  sync.RWMutex
}

type SecProductInfoConf struct {
	ProductID int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}

func InitConfig() (err error) {
	redisAddr, err := beego.AppConfig.String("redis_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_addr error:%v", err)
		return
	}
	etcdAddr, err := beego.AppConfig.String("etcd_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_addr error:%v", err)
		return
	}
	logs.Debug("read config succ, redis addr:%v", redisAddr)
	logs.Debug("read config succ, etcd addr:%v", etcdAddr)
	secKillConf.RedisAddr.RedisAddr = redisAddr
	secKillConf.EtcdConf.EtcdAddr = etcdAddr

	redisMaxIdle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_idle error:%v", err)
		return
	}
	redisMaxActive, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_active error:%v", err)
		return
	}
	redisIdleTimeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_idle_timeout error:%v", err)
		return
	}
	secKillConf.RedisAddr.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisAddr.RedisMaxActive = redisMaxActive
	secKillConf.RedisAddr.RedisIdleTimeout = redisIdleTimeout

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout err: %v", err)
		return
	}
	etcdSecKeyPrefix, err := beego.AppConfig.String("etcd_sec_key_prefix")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_sec_key_prefix error:%v", err)
		return
	}
	etcdProductKey, err := beego.AppConfig.String("etcd_product_key")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
		return
	}
	if !strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}
	secKillConf.EtcdConf.Timeout = etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix = etcdSecKeyPrefix
	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", etcdSecKeyPrefix, etcdProductKey)

	logPath, err := beego.AppConfig.String("log_path")
	if err != nil {
		err = fmt.Errorf("init config failed, read log_path err: %v", err)
		return
	}
	logLevel, err := beego.AppConfig.String("log_level")
	if err != nil {
		err = fmt.Errorf("init config failed, read log_path err: %v", err)
		return
	}
	secKillConf.LogPath = logPath
	secKillConf.LogLevel = logLevel

	return
}
