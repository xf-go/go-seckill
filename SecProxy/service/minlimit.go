package service

// SecLimit second limit
type MinLimit struct {
	count   int
	curTime int64
}

func (ml *MinLimit) Count(nowTime int64) (curCount int) {
	if ml.curTime != nowTime {
		ml.count = 1
		ml.curTime = nowTime
		curCount = ml.count
		return
	}

	ml.count++
	curCount = ml.count
	return
}

func (ml *MinLimit) Check(nowTime int64) int {
	if ml.curTime != nowTime {
		return 0
	}

	return ml.count
}
