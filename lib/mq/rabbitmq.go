/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-26 16:47:36
 * @LastEditTime: 2021-04-28 14:47:05
 */
package mq

import (
	"fmt"
	"iris-project/app/config"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ MQ实例
type RabbitMQ struct {
	uri     string
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect 连接服务器
func (mq *RabbitMQ) Connect() (err error) {
	mq.uri = fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		config.RabbitMQ.Username,
		config.RabbitMQ.Password,
		config.RabbitMQ.IP,
		config.RabbitMQ.Port,
	)
	mq.conn, err = amqp.Dial(mq.uri)
	if err != nil {
		log.Printf("[amqp] connect error: %s\n", err)
		return
	}
	return
}

// Channel 打开通道
func (mq *RabbitMQ) Channel() (err error) {
	mq.channel, err = mq.conn.Channel()
	if err != nil {
		log.Printf("[amqp] get channel error: %s\n", err)
		return
	}
	// 通道设置为Confirm模式
	err = mq.channel.Confirm(false)
	if err != nil {
		log.Printf("[amqp] set channel confirm error: %s\n", err)
		return
	}
	return
}
