package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
)

func initRedisPool(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConf.RedisAddr)
		},
	}
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed. err: %v", err)
		return
	}

	return
}

func initRedis(conf *SecLayerConf) (err error) {
	secLayerContext.proxy2LayerRedisPool, err = initRedisPool(conf.Proxy2LayerRedis)
	if err != nil {
		logs.Error("init proxy2layer redis failed. err:%v", err)
		return
	}

	secLayerContext.layer2ProxyRedisPool, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init layer2proxy redis failed. err:%v", err)
		return
	}

	return
}

func RunProcess() (err error) {
	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go handleRead()
	}

	for i := 0; i < secLayerContext.secLayerConf.WriteGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go handleWrite()
	}

	for i := 0; i < secLayerContext.secLayerConf.HandleUserGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go handleUser()
	}

	logs.Debug("all process goroutine started")
	secLayerContext.waitGroup.Wait()
	logs.Debug("wait all goroutine exited")
	return
}

func handleRead() {
	for {
		conn := secLayerContext.layer2ProxyRedisPool.Get()
		for {
			data, err := redis.Bytes(conn.Do("BLPOP", secLayerContext.secLayerConf.Proxy2LayerRedis.RedisQueueName, 0))
			if err != nil {
				logs.Error("pop from redis failed. err: %v", err)
				break
			}

			var req SecRequest
			err = json.Unmarshal(data, &req)
			if err != nil {
				logs.Error("json.Unmarshal failed. err: %v", err)
				continue
			}

			now := time.Now().Unix()
			if now-req.AccessTime.Unix() > secLayerContext.secLayerConf.MaxRequestWaitTimeout {
				logs.Warn("req is expired")
				continue
			}

			ticker := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToHandleChanTimeout))
			select {
			case secLayerContext.Read2HandleChan <- &req:
			case <-ticker.C:
				logs.Warn("send to handle chan timeout. req: %v", req)
			}
		}
		conn.Close()
	}
}

func handleWrite() {
	logs.Debug("handle write running.")
	for res := range secLayerContext.Handle2WriteChan {
		err := sendToRedis(res)
		if err != nil {
			logs.Error("send to redis failed. err: %v, res: %v", err, res)
			continue
		}
	}
}

func sendToRedis(res *SecResponse) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		logs.Error("json.Marshal failed. err: %v", err)
		return
	}

	conn := secLayerContext.layer2ProxyRedisPool.Get()
	_, err = conn.Do("RPUSH", secLayerContext.secLayerConf.Layer2ProxyRedis.RedisQueueName, string(data))
	if err != nil {
		logs.Error("rpush to redis failed. err: %v", err)
		return
	}

	return
}

func handleUser() {
	logs.Debug("handle user running.")
	for req := range secLayerContext.Read2HandleChan {
		logs.Debug("begin process request, req: ", req)
		res, err := handleSecKill(req)
		if err != nil {
			logs.Error("process request failed. err: %v, req: %v", err, req)
			res = &SecResponse{
				Code: ErrServiceBusy,
			}
		}

		ticker := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToWriteChanTimeout))
		select {
		case secLayerContext.Handle2WriteChan <- res:
		case <-ticker.C:
			logs.Warn("send to response chan timeout. req: %v, res: %v", req, res)
		}

	}
}

func handleSecKill(req *SecRequest) (res *SecResponse, err error) {
	secLayerContext.RWSecProductLock.RLock()
	defer secLayerContext.RWSecProductLock.RUnlock()

	res = &SecResponse{}
	res.UserId = req.UserId
	res.ProductId = req.ProductId
	product, ok := secLayerContext.secLayerConf.SecProductInfoMap[req.ProductId]
	if !ok {
		logs.Error("product[%d] not found", req.ProductId)
		res.Code = ErrNotFoundProduct
		return
	}

	if product.Status == ProductStatusSoldout {
		res.Code = ErrSoldout
		return
	}

	now := time.Now().Unix()
	alreadySoldout := product.secLimit.Check(now)
	if alreadySoldout >= product.soldMaxLimit {
		res.Code = ErrRetry
		return
	}

	// 限制每人购买数量
	secLayerContext.HistoryMapLock.Lock()
	userHistory, ok := secLayerContext.HistoryMap[req.UserId]
	if !ok {
		userHistory = &UserBuyHistory{
			history: make(map[int]int, 16),
		}

		secLayerContext.HistoryMap[req.UserId] = userHistory
	}

	historyCount := userHistory.GetProductBuyCount(req.ProductId)
	secLayerContext.HistoryMapLock.Unlock()

	if historyCount >= product.OnePersonBuyLimit {
		res.Code = ErrAlreadyBuy
		return
	}

	// 限制商品总数
	curSoldCount := secLayerContext.productCountMgr.Count(req.ProductId)
	if curSoldCount >= product.Total {
		res.Code = ErrSoldout
		product.Status = ProductStatusSoldout
		return
	}

	// 概率抽奖
	curRate := rand.Float64()
	if curRate > product.BuyRate {
		res.Code = ErrRetry
		return
	}

	// 更新总数
	userHistory.Add(req.ProductId, 1)
	secLayerContext.productCountMgr.Add(req.ProductId, 1)

	res.Code = ErrSecKillSucc
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId, req.ProductId, now, secLayerContext.secLayerConf.TokenPasswd)

	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = now

	return
}
