package model

import "github.com/streadway/amqp"

var MQ *amqp.Connection

func RabbitMQ(str string) {
	conn, err := amqp.Dial(str)
	if err != nil {
		panic(err)
	}
	MQ = conn
}
