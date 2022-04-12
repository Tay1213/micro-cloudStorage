package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"micro-cloudStorage/mq-server/conf"
	"strings"
	"time"
)

var DB *gorm.DB

func Setup() {
	pathDatabase := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.DatabaseSetting.User,
		conf.DatabaseSetting.Password,
		conf.DatabaseSetting.Host,
		conf.DatabaseSetting.Name)
	Database(pathDatabase)
	pathRabbitMQ := strings.Join([]string{conf.RabbitMqSetting.RabbitMQ, "://",
		conf.RabbitMqSetting.RabbitMQUser, ":", conf.RabbitMqSetting.RabbitMQPassWord, "@",
		conf.RabbitMqSetting.RabbitMQHost, ":", conf.RabbitMqSetting.RabbitMQPort, "/"}, "")
	RabbitMQ(pathRabbitMQ)
}

func Database(conn string) {
	db, err := gorm.Open(conf.DatabaseSetting.Type, conn)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	//默认不加复数
	db.SingularTable(true)
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(20)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db
	DB.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&FileSystem{})
}
