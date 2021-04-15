package service

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
)

func WriteHandle() {
	for {
		req := <-secKillServer.SecReqChan
		conn := secKillServer.proxy2LayerRedisPool.Get()

		data, err := json.Marshal(req)
		if err != nil {
			logs.Error("json.Marshal failed, err:%v req:%v", err, req)
			continue
		}

		_, err = conn.Do("LPUSH", "sec_queue", data)
		if err != nil {
			logs.Error("redis lpush failed, err:%v, req:%v", err, req)
			continue
		}

		conn.Close()
	}
	return
}

func ReadHandle() {
	return
}
