package model

import (
	"github.com/jmoiron/sqlx"
	etcd "go.etcd.io/etcd/client/v3"
)

var (
	DB             *sqlx.DB
	EtcdClient     *etcd.Client
	EtcdKeyPrefix  string
	EtcdProductKey string
)

func Init(db *sqlx.DB, etcdCli *etcd.Client, prefix, productKey string) {
	DB = db
	EtcdClient = etcdCli
	EtcdKeyPrefix = prefix
	EtcdProductKey = productKey
}
