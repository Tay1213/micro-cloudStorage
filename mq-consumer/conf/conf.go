package conf

import "github.com/go-ini/ini"

type Database struct {
	Type     string
	User     string
	Password string
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

	mapTo("database", DatabaseSetting)
	mapTo("rabbitmq", RabbitMqSetting)
}

func mapTo(str string, o interface{}) {
	err := conf.Section(str).MapTo(o)
	if err != nil {
		panic(err)
	}
}
