package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/clientv3"
)

const (
	EtcdKey = "/oldboy/backend/seckill/product"
)

type SecInfoConf struct {
	ProductID int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}

func main() {
	SetLogConfToEtcd()
}

func SetLogConfToEtcd() {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err: ", err)
		return
	}
	defer cli.Close()

	var secInfoConfArr []SecInfoConf
	secInfoConfArr = append(
		secInfoConfArr,
		SecInfoConf{
			ProductID: 1,
			StartTime: 1617962400,
			EndTime:   1617969600,
			Status:    0,
			Total:     1000,
			Left:      1000,
		},
	)
	secInfoConfArr = append(
		secInfoConfArr,
		SecInfoConf{
			ProductID: 2,
			StartTime: 1617962400,
			EndTime:   1617969600,
			Status:    0,
			Total:     2000,
			Left:      1000,
		},
	)

	data, err := json.Marshal(secInfoConfArr)
	if err != nil {
		fmt.Println("json.Marshal failed, err: ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err: ", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("%s : %s \n", kv.Key, kv.Value)
	}
}
