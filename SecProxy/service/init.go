package service

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
)

var (
	secKillServer *SecKillServer
)

func InitService(server *SecKillServer) (err error) {
	secKillServer = server

	err = loadBlackList()
	if err != nil {
		logs.Error("load black list failed, err:%v", err)
		return
	}

	logs.Debug("init service succ, config:%v", server)

	err = initProxy2LayerRedis()
	if err != nil {
		logs.Error("init proxy2Layer redis failed, err:%v", err)
		return
	}

	secKillServer.secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*Limit, 10000),
		IpLimitMap:   make(map[string]*Limit, 10000),
	}

	secKillServer.SecReqChan = make(chan *SecRequest, secKillServer.SecReqChanSize)

	initRedisProcessFunc()

	return
}

func initRedisProcessFunc() {
	for i := 0; i < secKillServer.WriteProxy2LayerGoroutineNum; i++ {
		go WriteHandle()
	}
	for i := 0; i < secKillServer.ReadLayer2ProxyGoroutineNum; i++ {
		go ReadHandle()
	}
}

func initProxy2LayerRedis() (err error) {
	secKillServer.proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:     secKillServer.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   secKillServer.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillServer.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillServer.RedisProxy2LayerConf.RedisAddr)
		},
	}
	conn := secKillServer.proxy2LayerRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err: %v", err)
		return
	}

	return
}

func initBlackRedis() (err error) {
	secKillServer.blackRedisPool = &redis.Pool{
		MaxIdle:     secKillServer.RedisBlackConf.RedisMaxIdle,
		MaxActive:   secKillServer.RedisBlackConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillServer.RedisBlackConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillServer.RedisBlackConf.RedisAddr)
		},
	}
	conn := secKillServer.blackRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err: %v", err)
		return
	}

	return
}

func loadBlackList() (err error) {
	err = initBlackRedis()
	if err != nil {
		logs.Error("init black redis failed, err:%v", err)
		return
	}

	conn := secKillServer.blackRedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", "idblacklist")
	idlist, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed,err:%v", err)
		return
	}

	for _, v := range idlist {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id [%v]", id)
			continue
		}
		secKillServer.IdBlackMap[id] = true
	}

	reply, err = conn.Do("hgetall", "ipblacklist")
	iplist, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}

	for _, v := range iplist {
		secKillServer.IpBlackMap[v] = true
	}

	go syncIpBlackList()
	go syncIdBlackList()

	return
}

func syncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()
	for {
		conn := secKillServer.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackiplist", time.Second)
		ip, err := redis.String(reply, err)
		if err != nil {
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 100 || curTime-lastTime > 5 {
			secKillServer.RWBlackLock.Lock()
			for _, v := range ipList {
				secKillServer.IpBlackMap[v] = true
			}
			secKillServer.RWBlackLock.Unlock()

			lastTime = curTime
			logs.Info("sync ip list from redis succ, ip[%v]", ipList)
		}
	}
}

func syncIdBlackList() {
	for {
		conn := secKillServer.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackidlist", time.Second)
		id, err := redis.Int(reply, err)
		if err != nil {
			continue
		}
		secKillServer.RWBlackLock.Lock()
		secKillServer.IdBlackMap[id] = true
		secKillServer.RWBlackLock.Unlock()

		logs.Info("sync id list from redis succ, id[%v]", id)
	}
}
