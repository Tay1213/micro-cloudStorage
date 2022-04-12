package conf

import "github.com/go-ini/ini"

type Database struct {
	User     string
	Password string
	Type     string
	Host     string
	Name     string
}

type Rabbit struct {
	RabbitMQ         string
	RabbitMQUser     string
	RabbitMQPassWord string
	RabbitMQHost     string
	RabbitMQPort     string
}

var DatabaseSetting = &Database{}
var RabbitMqSetting = &Rabbit{}
var conf *ini.File

func Setup() {
	var err error
	conf, err = ini.Load("conf/conf.ini")
	if err != nil {
		panic(err)
	}
	MapTo("database", DatabaseSetting)
	MapTo("rabbitmq", RabbitMqSetting)
}

func MapTo(str string, o interface{}) {
	err := conf.Section(str).MapTo(o)
	if err != nil {
		panic(err)
	}
}
