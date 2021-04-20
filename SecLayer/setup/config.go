package setup

import (
	"SecLayer/service"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/adapter/config"
	"github.com/beego/beego/v2/core/logs"
)

var AppConfig *service.SecLayerConf

func InitConfig(adapterName, filename string) (err error) {
	conf, err := config.NewConfig(adapterName, filename)
	if err != nil {
		logs.Error("init config failed. err:%v", err)
		return
	}

	// 读取日志库配置
	AppConfig = &service.SecLayerConf{}
	AppConfig.LogLevel = conf.String("logs::log_level")
	if len(AppConfig.LogLevel) == 0 {
		AppConfig.LogLevel = "debug"
	}

	AppConfig.LogPath = conf.String("logs::log_path")
	if len(AppConfig.LogPath) == 0 {
		AppConfig.LogPath = "./logs"
	}

	// 读取redis相关配置
	redisAddr := conf.String("redis::redis_proxy2layer_addr")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_addr failed. err:%v", err)
		return
	}

	redisMaxIdle, err := conf.Int("redis::redis_proxy2layer_max_idle")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_max_idle failed. err:%v", err)
		return
	}

	redisMaxActive, err := conf.Int("redis::redis_proxy2layer_max_active")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_max_active failed. err:%v", err)
		return
	}

	redisIdleTimeout, err := conf.Int("redis::redis_proxy2layer_idle_timeout")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_idle_timeout failed. err:%v", err)
		return
	}

	redisQueueName := conf.String("redis::redis_proxy2layer_queue_name")
	if len(redisQueueName) == 0 {
		logs.Error("read redis::redis_proxy2layer_queue_name failed. err:%v", err)
		return
	}

	AppConfig.Proxy2LayerRedis.RedisAddr = redisAddr
	AppConfig.Proxy2LayerRedis.RedisMaxActive = redisMaxActive
	AppConfig.Proxy2LayerRedis.RedisMaxIdle = redisMaxIdle
	AppConfig.Proxy2LayerRedis.RedisIdleTimeout = redisIdleTimeout
	AppConfig.Proxy2LayerRedis.RedisQueueName = redisQueueName

	// 读取redis相关配置
	redisAddr = conf.String("redis::redis_proxy2layer_addr")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_addr failed. err:%v", err)
		return
	}

	redisMaxIdle, err = conf.Int("redis::redis_proxy2layer_max_idle")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_max_idle failed. err:%v", err)
		return
	}

	redisMaxActive, err = conf.Int("redis::redis_proxy2layer_max_active")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_max_active failed. err:%v", err)
		return
	}

	redisIdleTimeout, err = conf.Int("redis::redis_proxy2layer_idle_timeout")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_idle_timeout failed. err:%v", err)
		return
	}

	redisQueueName = conf.String("redis::redis_layer2proxy_queue_name")
	if len(redisQueueName) == 0 {
		logs.Error("read redis::redis_layer2proxy_queue_name failed. err:%v", err)
		return
	}

	AppConfig.Layer2ProxyRedis.RedisAddr = redisAddr
	AppConfig.Layer2ProxyRedis.RedisMaxActive = redisMaxActive
	AppConfig.Layer2ProxyRedis.RedisMaxIdle = redisMaxIdle
	AppConfig.Layer2ProxyRedis.RedisIdleTimeout = redisIdleTimeout
	AppConfig.Layer2ProxyRedis.RedisQueueName = redisQueueName

	// etcd
	etcdAddr := conf.String("etcd::etcd_addr")
	if len(etcdAddr) == 0 {
		err = fmt.Errorf("init config failed, read etcd_addr error")
		return
	}
	etcdTimeout, err := conf.Int("etcd::etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout err: %v", err)
		return
	}
	etcdSecKeyPrefix := conf.String("etcd::etcd_sec_key_prefix")
	if len(etcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_key_prefix error: %v", err)
		return
	}
	etcdSecProductKey := conf.String("etcd::etcd_sec_product_key")
	if len(etcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_product_key error: %v", err)
		return
	}
	AppConfig.EtcdConfig.EtcdAddr = etcdAddr
	AppConfig.EtcdConfig.Timeout = etcdTimeout
	AppConfig.EtcdConfig.EtcdSecKeyPrefix = etcdSecKeyPrefix
	if !strings.HasSuffix(AppConfig.EtcdConfig.EtcdSecKeyPrefix, "/") {
		AppConfig.EtcdConfig.EtcdSecKeyPrefix = AppConfig.EtcdConfig.EtcdSecKeyPrefix + "/"
	}
	AppConfig.EtcdConfig.EtcdSecProductKey = fmt.Sprintf("%s%s", AppConfig.EtcdConfig.EtcdSecKeyPrefix, etcdSecProductKey)

	writeGoroutineNum, err := conf.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		logs.Error("read service::write_proxy2layer_goroutine_num failed. err:%v", err)
		return
	}

	readGoroutineNum, err := conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		logs.Error("read service::read_layer2proxy_goroutine_num failed. err:%v", err)
		return
	}

	handleUserGoroutineNum, err := conf.Int("service::handle_user_goroutine_num")
	if err != nil {
		logs.Error("read service::read_layer2proxy_goroutine_num failed. err:%v", err)
		return
	}

	read2handleChanSize, err := conf.Int("service::read2handle_chan_size")
	if err != nil {
		logs.Error("read service::read2handle_chan_size failed. err:%v", err)
		return
	}

	handle2writeChanSize, err := conf.Int("service::handle2write_chan_size")
	if err != nil {
		logs.Error("read service::handle2write_chan_size failed. err:%v", err)
		return
	}

	maxRequestWaitTimeout, err := conf.Int64("service::max_request_wait_timeout")
	if err != nil {
		logs.Error("read service::max_request_wait_timeout failed. err:%v", err)
		return
	}

	sendToWriteChanTimeout, err := conf.Int64("service::send_to_write_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_write_chan_timeout failed. err:%v", err)
		return
	}

	sendToHandleChanTimeout, err := conf.Int64("service::send_to_handle_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_handle_chan_timeout failed. err:%v", err)
		return
	}

	AppConfig.WriteGoroutineNum = writeGoroutineNum
	AppConfig.ReadGoroutineNum = readGoroutineNum
	AppConfig.HandleUserGoroutineNum = handleUserGoroutineNum
	AppConfig.Read2HandleChanSize = read2handleChanSize
	AppConfig.Handle2WriteChanSize = handle2writeChanSize
	AppConfig.MaxRequestWaitTimeout = maxRequestWaitTimeout
	AppConfig.SendToWriteChanTimeout = sendToWriteChanTimeout
	AppConfig.SendToHandleChanTimeout = sendToHandleChanTimeout

	tokenPasswd := conf.String("service::seckill_token_passwd")
	if len(tokenPasswd) == 0 {
		logs.Error("read service::seckill_token_passwd failed. err:%v", err)
		return
	}
	AppConfig.TokenPasswd = tokenPasswd

	return
}
