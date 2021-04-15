package setup

import (
	"SecProxy/service"
	"fmt"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

var (
	secKillServer = &service.SecKillServer{
		SecProductInfoMap: make(map[int]*service.SecProductInfoConf, 1024),
	}
)

func InitConfig() (err error) {
	// blacklist redis
	redisAddr, err := beego.AppConfig.String("redis_black_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_addr error:%v", err)
		return
	}
	redisMaxIdle, err := beego.AppConfig.Int("redis_black_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_max_idle error:%v", err)
		return
	}
	redisMaxActive, err := beego.AppConfig.Int("redis_black_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_max_active error:%v", err)
		return
	}
	redisIdleTimeout, err := beego.AppConfig.Int("redis_black_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_idle_timeout error:%v", err)
		return
	}
	secKillServer.RedisBlackConf.RedisAddr = redisAddr
	secKillServer.RedisBlackConf.RedisMaxIdle = redisMaxIdle
	secKillServer.RedisBlackConf.RedisMaxActive = redisMaxActive
	secKillServer.RedisBlackConf.RedisIdleTimeout = redisIdleTimeout

	// proxy2layer redis
	redisAddr, err = beego.AppConfig.String("redis_proxy2layer_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_addr error:%v", err)
		return
	}
	redisMaxIdle, err = beego.AppConfig.Int("redis_proxy2layer_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_max_idle error:%v", err)
		return
	}
	redisMaxActive, err = beego.AppConfig.Int("redis_proxy2layer_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_max_active error:%v", err)
		return
	}
	redisIdleTimeout, err = beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_idle_timeout error:%v", err)
		return
	}
	secKillServer.RedisProxy2LayerConf.RedisAddr = redisAddr
	secKillServer.RedisProxy2LayerConf.RedisMaxIdle = redisMaxIdle
	secKillServer.RedisProxy2LayerConf.RedisMaxActive = redisMaxActive
	secKillServer.RedisProxy2LayerConf.RedisIdleTimeout = redisIdleTimeout

	// etcd
	etcdAddr, err := beego.AppConfig.String("etcd_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_addr error:%v", err)
		return
	}
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
	etcdSecProductKey, err := beego.AppConfig.String("etcd_sec_product_key")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_sec_product_key error:%v", err)
		return
	}
	secKillServer.EtcdConf.EtcdAddr = etcdAddr
	secKillServer.EtcdConf.Timeout = etcdTimeout
	secKillServer.EtcdConf.EtcdSecKeyPrefix = etcdSecKeyPrefix
	if !strings.HasSuffix(secKillServer.EtcdConf.EtcdSecKeyPrefix, "/") {
		secKillServer.EtcdConf.EtcdSecKeyPrefix = secKillServer.EtcdConf.EtcdSecKeyPrefix + "/"
	}
	secKillServer.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillServer.EtcdConf.EtcdSecKeyPrefix, etcdSecProductKey)

	// log
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
	secKillServer.LogPath = logPath
	secKillServer.LogLevel = logLevel

	// cookie
	secretkey, err := beego.AppConfig.String("cookie_secretkey")
	if err != nil {
		err = fmt.Errorf("init config failed, read cookie_secretkey err: %v", err)
		return
	}
	secKillServer.CookieSecretKey = secretkey

	// access limit
	userSecAccessLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_sec_access_limit err: %v", err)
		return
	}
	ipSecAccessLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit err: %v", err)
		return
	}
	userMinAccessLimit, err := beego.AppConfig.Int("user_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_min_access_limit err: %v", err)
		return
	}
	ipMinAccessLimit, err := beego.AppConfig.Int("ip_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_min_access_limit err: %v", err)
		return
	}
	secKillServer.AccessLimitConf.UserSecAccessLimit = userSecAccessLimit
	secKillServer.AccessLimitConf.IPSecAccessLimit = ipSecAccessLimit
	secKillServer.AccessLimitConf.UserMinAccessLimit = userMinAccessLimit
	secKillServer.AccessLimitConf.IPMinAccessLimit = ipMinAccessLimit

	//
	refererWhitelist, err := beego.AppConfig.String("referer_whitelist")
	if err != nil {
		err = fmt.Errorf("init config failed, read referer_whitelist err: %v", err)
		return
	}
	if len(refererWhitelist) > 0 {
		secKillServer.RefererWhiteList = strings.Split(refererWhitelist, ",")
	}

	// write_proxy2layer_goroutine_num
	writeProxy2LayerGoroutineNum, err := beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read write_proxy2layer_goroutine_num err: %v", err)
		return
	}
	secKillServer.WriteProxy2LayerGoroutineNum = writeProxy2LayerGoroutineNum

	// write_proxy2layer_goroutine_num
	readLayer2ProxyGoroutineNum, err := beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read read_layer2proxy_goroutine_num err: %v", err)
		return
	}
	secKillServer.ReadLayer2ProxyGoroutineNum = readLayer2ProxyGoroutineNum

	return
}
