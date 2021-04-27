package service

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

func SecKill(req *SecRequest) (data map[string]interface{}, code int, err error) {
	secKillServer.RWSecProductLock.RLock()
	defer secKillServer.RWSecProductLock.RUnlock()

	err = userCheck(req)
	if err != nil {
		code = ErrUserCheckAuthFailed
		logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return
	}

	err = antiSpam(req)
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return
	}

	data, code, err = SecInfoById(req.ProductId)
	if err != nil {
		logs.Warn("userId[%d] SecInfoById failed, req[%v]", req.UserId, req)
		return
	}
	if code != 0 {
		logs.Warn("userId[%d] SecInfoById failed, code[%d] req[%v]", req.UserId, code, req)
		return
	}
	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)

	secKillServer.SecReqChan <- req

	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		secKillServer.UserConnMapLock.Lock()
		delete(secKillServer.UserConnMap, userKey)
		secKillServer.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		return nil, ErrProcessTimeout, fmt.Errorf("request timeout")
	case <-req.CloseNotify:
		return nil, ErrClientClosed, fmt.Errorf("client already closed")
	case result := <-req.ResultChan:
		if result.Code != 1002 {
			return data, code, fmt.Errorf("client already closed")
		}
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return data, code, nil
	}

}

func userCheck(req *SecRequest) (err error) {
	found := false
	for _, referer := range secKillServer.RefererWhiteList {
		if referer == req.ClientReferer {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("invalid request")
		logs.Warn("user[%d] is rejected by referer, req[%v]", req.UserId, req)
		return
	}
	authData := fmt.Sprintf("%d:%s", req.UserId, secKillServer.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))
	if authSign != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}
	return
}

func SecInfoList() (data []map[string]interface{}, code int, err error) {
	secKillServer.RWSecProductLock.RLock()
	defer secKillServer.RWSecProductLock.RUnlock()

	for _, v := range secKillServer.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			logs.Error("get product_id[%d] failed, err:%v", v.ProductId, err)
			continue
		}
		data = append(data, item)
	}
	return
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {
	secKillServer.RWSecProductLock.RLock()
	defer secKillServer.RWSecProductLock.RUnlock()

	item, code, err := SecInfoById(productId)
	if err != nil {
		return
	}

	data = append(data, item)
	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {
	secKillServer.RWSecProductLock.RLock()
	defer secKillServer.RWSecProductLock.RUnlock()

	v, ok := secKillServer.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	start, end := false, false
	status := "success"

	now := time.Now().Unix()
	if now < v.StartTime {
		status = "sec kill is not start"
		code = ErrActivityNotStart
	}
	if now > v.StartTime {
		start = true
	}
	if now > v.EndTime {
		end = true
		status = "sec kill is alredy end"
		code = ErrActivityAlreadyEnd
	}
	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		status = "product is sale out"
		code = ErrActivitySaleOut
	}
	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = status

	return
}
