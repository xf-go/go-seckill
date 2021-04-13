package service

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
)

var (
	secKillConf *SecKillConf
)

func InitService(conf *SecKillConf) {
	secKillConf = conf
	logs.Debug("init service succ, config:%v", conf)
}

func SecInfo(productId int) (data map[string]interface{}, code int, err error) {
	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.RUnlock()

	v, ok := secKillConf.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status
	return
}
