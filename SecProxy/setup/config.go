package setup

import (
	"SecProxy/service"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

var (
	secKillConf = &service.SecKillConf{
		SecProductInfoMap: make(map[int]*service.SecProductInfoConf, 1024),
	}
)

func InitConfig() (err error) {
	redisAddr, err := beego.AppConfig.String("redis_black_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_addr error:%v", err)
		return
	}
	etcdAddr, err := beego.AppConfig.String("etcd_addr")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_addr error:%v", err)
		return
	}
	logs.Debug("read config succ, redis addr:%v", redisAddr)
	logs.Debug("read config succ, etcd addr:%v", etcdAddr)
	secKillConf.RedisBlackAddr.RedisAddr = redisAddr
	secKillConf.EtcdConf.EtcdAddr = etcdAddr

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
	secKillConf.RedisBlackAddr.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisBlackAddr.RedisMaxActive = redisMaxActive
	secKillConf.RedisBlackAddr.RedisIdleTimeout = redisIdleTimeout

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

	secKillConf.EtcdConf.Timeout = etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix = etcdSecKeyPrefix
	if !strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}
	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillConf.EtcdConf.EtcdSecKeyPrefix, etcdSecProductKey)
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

	secretkey, err := beego.AppConfig.String("cookie_secretkey")
	if err != nil {
		err = fmt.Errorf("init config failed, read cookie_secretkey err: %v", err)
		return
	}
	secKillConf.CookieSecretKey = secretkey

	secLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_sec_access_limit err: %v", err)
		return
	}
	secKillConf.UserSecAccessLimit = secLimit

	refererWhitelist, err := beego.AppConfig.String("referer_whitelist")
	if err != nil {
		err = fmt.Errorf("init config failed, read referer_whitelist err: %v", err)
		return
	}
	if len(refererWhitelist) > 0 {
		secKillConf.RefererWhiteList = strings.Split(refererWhitelist, ",")
	}

	ipSecAccessLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit err: %v", err)
		return
	}
	secKillConf.IPSecAccessLimit = ipSecAccessLimit

	return
}
