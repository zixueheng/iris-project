package global

import (
	"fmt"
	"iris-project/app/config"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var MqConn *amqp.Connection

func init() {
	if !config.RabbitMQ.On {
		return
	}
	var (
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			config.RabbitMQ.Username,
			config.RabbitMQ.Password,
			config.RabbitMQ.IP,
			config.RabbitMQ.Port,
		)
		err error
	)

	// 建立amqp链接
	MqConn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Println(err.Error())
		for {
			time.Sleep(60 * time.Second) // 等待60秒再次连接
			log.Println("再次连接rabbitMQ")
			MqConn, err = amqp.Dial(rabbitURL)
			if err == nil {
				break
			} else {
				log.Println(err.Error())
			}
		}
	}
}
