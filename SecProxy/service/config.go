package service

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
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
	RedisBlackAddr     RedisConf
	EtcdConf           EtcdConf
	LogPath            string
	LogLevel           string
	SecProductInfoMap  map[int]*SecProductInfoConf
	RWSecProductLock   sync.RWMutex
	CookieSecretKey    string
	UserSecAccessLimit int
	RefererWhiteList   []string
	IPSecAccessLimit   int
	IpBlackList        map[string]bool
	IdBlackList        map[int]bool

	BlackRedisPool *redis.Pool
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientReferer string
}
