package main

import (
	"micro-cloudStorage/mq-server/conf"
	"micro-cloudStorage/mq-server/model"
	"micro-cloudStorage/mq-server/service"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	conf.Setup()
	model.Setup()

	forever := make(chan bool)
	service.CreateFile()
	<-forever
}
