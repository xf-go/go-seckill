package service

import (
	"fmt"
	"sync"
)

var (
	secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*SecLimit, 10000),
	}
)

type SecLimitMgr struct {
	UserLimitMap map[int]*SecLimit
	IpLimitMap   map[string]*SecLimit
	lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	secLimitMgr.lock.Lock()

	userLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		userLimit = &SecLimit{}
		secLimitMgr.UserLimitMap[req.UserId] = userLimit
	}
	count := userLimit.Count(req.AccessTime.Unix())
	if count > secKillServer.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	ipLimit, ok := secLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		ipLimit = &SecLimit{}
		secLimitMgr.IpLimitMap[req.ClientAddr] = ipLimit
	}
	count = ipLimit.Count(req.AccessTime.Unix())
	if count > secKillServer.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	secLimitMgr.lock.Unlock()

	return
}

// SecLimit second limit
type SecLimit struct {
	count   int
	curTime int64
}

func (s *SecLimit) Count(nowTime int64) (curCount int) {
	if s.curTime != nowTime {
		s.count = 1
		s.curTime = nowTime
		curCount = s.count
		return
	}

	s.count++
	curCount = s.count
	return
}

func (s *SecLimit) Check(nowTime int64) int {
	if s.curTime != nowTime {
		return 0
	}

	return s.count
}
