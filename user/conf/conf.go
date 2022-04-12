package conf

import (
	"github.com/go-ini/ini"
	"time"
)

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var conf *ini.File

func Setup() {
	var err error
	conf, err = ini.Load("conf/conf.ini")
	if err != nil {

	}
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func mapTo(str string, s interface{}) {
	err := conf.Section(str).MapTo(s)
	if err != nil {

	}
}
