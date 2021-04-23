package service

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
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

		_, err = conn.Do("LPUSH", "sec_queue", string(data))
		if err != nil {
			logs.Error("redis lpush failed. err:%v, req:%v", err, req)
			conn.Close()
			continue
		}

		conn.Close()
	}
}

func ReadHandle() {
	for {
		conn := secKillServer.proxy2LayerRedisPool.Get()

		reply, err := conn.Do("BRPOP", "recv_queue", 0)
		data, err := redis.String(reply, err)
		// if err == redis.ErrNil {
		// 	conn.Close()
		// 	time.Sleep(time.Second)
		// 	continue
		// }
		if err != nil {
			logs.Error("redis rpop failed. err:%v", err)
			conn.Close()
			continue
		}

		var result SecResult
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			logs.Error("json.Unmarshal failed. err:%v", err)
			conn.Close()
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		secKillServer.UserConnMapLock.Lock()
		resultChan, ok := secKillServer.UserConnMap[userKey]
		secKillServer.UserConnMapLock.Unlock()
		if !ok {
			conn.Close()
			logs.Error("user not found. err: %v", err)
			continue
		}

		resultChan <- &result
		conn.Close()
	}
}
