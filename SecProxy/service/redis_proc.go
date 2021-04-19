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
			logs.Error("json.Marshal failed. err:%v, req:%v", err, req)
			conn.Close()
			continue
		}

		_, err = conn.Do("LPUSH", "sec_queue", data)
		if err != nil {
			logs.Error("redis lpush failed. err:%v, req:%v", err, req)
			conn.Close()
			continue
		}

		conn.Close()
	}
}

func ReadHandle() {
	return
	// for {
	// 	conn := secKillServer.proxy2LayerRedisPool.Get()

	// 	data, err := conn.Do("RPOP", "sec_queue")
	// 	if err != nil {
	// 		logs.Error("redis rpop failed. err:%v", err)
	// 		conn.Close()
	// 		continue
	// 	}

	// 	var ch secKillServer.SecReqChan
	// 	json.Unmarshal([]byte(data), &ch)
	// }
}
