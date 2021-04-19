package service

import (
	"fmt"
	"sync"

	"github.com/beego/beego/v2/core/logs"
)

var secLimitMgr = &SecLimitMgr{
	UserLimitMap: make(map[int]*Limit, 10000),
	IpLimitMap:   make(map[string]*Limit, 10000),
}

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	// 判断用户Id是否在黑名单
	_, ok := secKillServer.IdBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("user[%d] is block by id black list", req.UserId)
		return
	}

	// 判断客户端IP是否在黑名单
	_, ok = secKillServer.IpBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("user[%d] ip[%s] is block by id black list", req.UserId, req.ClientAddr)
		return
	}

	var secIdCount, minIdCount, secIpCount, minIpCount int
	// 加锁
	secLimitMgr.lock.Lock()
	{
		// 用户Id频率控制
		limit, ok := secLimitMgr.UserLimitMap[req.UserId]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			secLimitMgr.UserLimitMap[req.UserId] = limit
		}

		secIdCount = limit.secLimit.Count(req.AccessTime.Unix())
		minIdCount = limit.minLimit.Count(req.AccessTime.Unix())

		// 客户端Ip频率控制
		limit, ok = secLimitMgr.IpLimitMap[req.ClientAddr]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			secLimitMgr.IpLimitMap[req.ClientAddr] = limit
		}

		secIpCount = limit.secLimit.Count(req.AccessTime.Unix())
		minIpCount = limit.minLimit.Count(req.AccessTime.Unix())
	}
	// 释放锁
	secLimitMgr.lock.Unlock()

	// 判断该用户一秒内访问次数是否大于配置的最大访问次数
	if secIdCount > secKillServer.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}
	// 判断该用户一分钟内访问次数是否大于配置的最大访问次数
	if minIdCount > secKillServer.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}
	// 判断该IP一秒内访问次数是否大于配置的最大访问次数
	if secIpCount > secKillServer.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}
	// 判断该IP一分钟内访问次数是否大于配置的最大访问次数
	if minIpCount > secKillServer.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	return
}

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(nowTime int64) int
}

type Limit struct {
	secLimit TimeLimit
	minLimit TimeLimit
}
