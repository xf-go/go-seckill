package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ID        int     `db:"id"`
	Name      string  `db:"name"`
	ProductId int     `db:"product_id"`
	StartTime int64   `db:"start_time"`
	EndTime   int64   `db:"end_time"`
	Total     int64   `db:"total"`
	Status    int8    `db:"status"`
	Speed     int     `db:"sec_speed"`
	BuyLimit  int     `db:"buy_limit"`
	BuyRate   float64 `db:"buy_rate"`

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
}

type SecProductInfoConf struct {
	ProductId         int
	StartTime         int64
	EndTime           int64
	Status            int8
	Total             int64
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	// 每秒最多能卖多少个
	soldMaxLimit int
	// 限速控制
	// secLimit *SecLimit
}

type ActivityModel struct {
}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (p *ActivityModel) GetActivityList() (activityList []*Activity, err error) {
	sql := "select id, name, product_id, start_time, end_time, total, status, sec_speed, buy_limit, buy_rate from activity order by od desc"
	err = DB.Select(&activityList, sql)
	if err != nil {
		logs.Error("select activity from database failed. err: %v", err)
		return
	}

	for _, v := range activityList {
		t := time.Unix(v.StartTime, 0)
		v.StartTimeStr = t.Format("2006-01-02 15:04:05")

		t = time.Unix(v.EndTime, 0)
		v.EndTimeStr = t.Format("2006-01-02 15:04:05")

		now := time.Now().Unix()
		if now > v.EndTime {
			v.StatusStr = "已结束"
			continue
		}

		if v.Status == ActivityStatusNormal {
			v.StatusStr = "正常"
		} else if v.Status == ActivityStatusDisable {
			v.StatusStr = "已禁用"
		}
	}

	return
}

func (p *ActivityModel) ProductValid(productId int, total int64) (valid bool, err error) {
	sql := "select id, total from product where product_id=?"
	var product Product
	err = DB.Select(&product, sql)
	if err != nil {
		logs.Error("select product from database failed, err: %v, sql: %v", err, sql)
		return
	}
	if product.ID <= 0 {
		err = fmt.Errorf("product id[%d] doen not exist", productId)
		return
	}
	if total > product.Total {
		err = fmt.Errorf("invalid product total[%d]", product.Total)
		return
	}

	valid = true
	return
}

func (p *ActivityModel) CreateActivity(activity *Activity) (err error) {
	valid, err := p.ProductValid(activity.ProductId, activity.Total)
	if err != nil {
		logs.Error("product valid failed. err: %v", err)
		return
	}
	if !valid {
		err = fmt.Errorf("invalid product id[%d]", activity.ProductId)
		logs.Error(err)
		return
	}

	if activity.StartTime <= 0 || activity.EndTime <= 0 {
		err = fmt.Errorf("invalid start[%v]|end[%v] time", activity.StartTime, activity.EndTime)
		logs.Error(err)
		return
	}
	if activity.EndTime <= activity.StartTime {
		err = fmt.Errorf("start[%v] time is greater than end[%v] time", activity.StartTime, activity.EndTime)
		logs.Error(err)
		return
	}
	now := time.Now().Unix()
	if activity.EndTime <= now || activity.StartTime <= now {
		err = fmt.Errorf("start[%v]|end[%v] time is less than now[%v]", activity.StartTime, activity.EndTime, now)
		logs.Error(err)
		return
	}

	sql := "insert into activity(name, product_id, start_time, end_time, total, sec_speed, buy_limit, buy_rate) values (?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(sql, activity.Name, activity.ProductId, activity.StartTime, activity.EndTime, activity.Total, activity.Speed, activity.BuyLimit, activity.BuyRate)
	if err != nil {
		logs.Error("create activity failed, err: %v, sql: %v", err, sql)
		return
	}
	logs.Debug("insert into database succ.")

	err = p.SyncToEtcd(activity)
	if err != nil {
		logs.Error("sync to etcd failed. err: %v, data[%v]", err, activity)
		return
	}

	return
}

func (p *ActivityModel) SyncToEtcd(activity *Activity) (err error) {
	if !strings.HasSuffix(EtcdKeyPrefix, "/") {
		EtcdKeyPrefix = EtcdKeyPrefix + "/"
	}
	etcdKey := fmt.Sprintf("%s%s", EtcdKeyPrefix, EtcdProductKey)
	secProductInfoList, err := loadProductFromEtcd(etcdKey)
	if err != nil {
		logs.Error("load product from etcd failed. err: %v", err)
		return
	}

	var secProductInfo SecProductInfoConf
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.EndTime = activity.EndTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.BuyRate = activity.BuyRate
	secProductInfo.soldMaxLimit = activity.Speed

	secProductInfoList = append(secProductInfoList, secProductInfo)

	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		logs.Error("json.Marshal failed. err: %v", err)
		return
	}

	EtcdClient.Put(context.Background(), etcdKey, string(data))
	if err != nil {
		logs.Error("put to etcd failed. err: %v, data[%v]", err, string(data))
		return
	}

	return
}

func loadProductFromEtcd(key string) (secProductInfo []SecProductInfoConf, err error) {
	logs.Debug("start load product from etcd.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := EtcdClient.Get(ctx, key)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err: %v", key, err)
		return
	}
	logs.Debug("load product from etcd srcc. resp: %v", resp)

	for k, v := range resp.Kvs {
		logs.Debug("key[%v] value[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("json.Unmarshal failed, err: %v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	// updateSecProductInfo(conf, secProductInfo)
	// initSecProductWatcher(conf)

	return
}
