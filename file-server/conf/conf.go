package conf

import (
	"github.com/go-ini/ini"
	"time"
)

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var conf *ini.File

var RedisSetting = &Redis{}

func Setup() {
	var err error
	conf, err = ini.Load("conf/conf.ini")
	if err != nil {
		panic(err)
	}
	mapTo("redis", RedisSetting)
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func mapTo(str string, o interface{}) {
	err := conf.Section(str).MapTo(o)
	if err != nil {
		panic(err)
	}
}
