/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-09 15:19:31
 * @LastEditTime: 2023-02-09 16:50:54
 */
package main

import (
	"fmt"
	mq "iris-project/lib/mq2"
	"log"
	"time"
)

// 生产者
func main() {
	rabbit := mq.NewRabbitMQ("yoyo_exchange", "direct", "yoyo_route", "yoyo_queue")
	defer rabbit.Close()

	for i := 1; i <= 3; i++ {
		go func(i int) {
			rabbit.SendMessage(mq.Message{Body: fmt.Sprintf("普通消息%d", i)})
			rabbit.SendDelayMessage(mq.Message{Body: fmt.Sprintf("延时5秒的消息%d", i), DelayTime: 5})
		}(i)
	}

	time.Sleep(time.Second * 3)
	log.Println("3秒后")

	rabbit.SendMessage(mq.Message{Body: "普通消息4"})
	rabbit.SendDelayMessage(mq.Message{Body: "延时5秒的消息4", DelayTime: 5})

}
