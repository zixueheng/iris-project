/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-09 14:25:30
 * @LastEditTime: 2023-02-10 10:17:37
 */
package main

import (
	"fmt"
	mq "iris-project/lib/mq2"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 消费者，先新启动
func main() {
	forever := make(chan bool)

	for i := 1; i <= 3; i++ {
		go func(i int) {
			// 新建连接
			rabbit := mq.NewRabbitMQ("yoyo_exchange", "direct", "yoyo_route", "yoyo_queue")
			// 一般来说消费者不关闭，常驻进程进行消息消费处理
			// defer rabbit.Close()

			// 执行消费
			rabbit.Consume(func(d amqp.Delivery) {
				//logger.Info("rabbitmq", zap.String("rabbitmq", string(d.Body)))
				log.Println(fmt.Sprintf("消费者%d收到消息：", i), string(d.Body))
			})
		}(i)
	}

	<-forever
}
