package main

var (
	secKillConf = &SecKillConf{}
)

type SecKillConf struct {
	RedisAddr string
	EtcdAddr  string
}
