package service

// SecLimit second limit
type SecLimit struct {
	count   int
	curTime int64
}

func (sl *SecLimit) Count(nowTime int64) (curCount int) {
	if sl.curTime != nowTime {
		sl.count = 1
		sl.curTime = nowTime
		curCount = sl.count
		return
	}

	sl.count++
	curCount = sl.count
	return
}

func (sl *SecLimit) Check(nowTime int64) int {
	if sl.curTime != nowTime {
		return 0
	}

	return sl.count
}
