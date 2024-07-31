/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2024-07-31 10:16:10
 * @LastEditTime: 2024-07-31 15:42:57
 */
package mq

import (
	"context"
	"iris-project/global"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetConn() *amqp.Connection {
	return global.MqConn
}

// CloseConn 关闭连接
func CloseConn() {
	GetConn().Close()
}

// CloseChannel 打开一个新通道
func OpenChannel() (*amqp.Channel, error) {
	return GetConn().Channel()
}

// CloseChannel 关闭通道
func CloseChannel(channel *amqp.Channel) {
	channel.Close()
}

// 交换机类型
const (
	DirectType  = "direct"
	FanoutType  = "fanout"
	TopicType   = "topic"
	HeadersType = "headers"
)

type (
	// 消息体：DelayTime 仅在 SendDelayMessage 方法有效
	Message struct {
		DelaySecond int // desc:延迟时间(秒)
		Body        string
	}

	MessageQueue struct {
		ExchangeName string // 交换器名称
		ExchangeKind string // 交换器类型 direct, fanout, topic, headers
		RouteKey     string // 路由名称
		QueueName    string // 队列名称
	}

	// 消费者回调方法
	Consumer func(amqp.Delivery)
)

func New(exchangeName, exchangeKind, routeKey, queueName string) (*MessageQueue, error) {
	var mq = &MessageQueue{
		ExchangeName: exchangeName,
		ExchangeKind: exchangeKind,
		RouteKey:     routeKey,
		QueueName:    queueName,
	}

	channel, err := OpenChannel()
	if err != nil {
		log.Printf("打开通道失败：%s\n", err.Error())
		return nil, err
	}
	defer CloseChannel(channel)

	if err := mq.declareExchange(channel, exchangeName, exchangeKind, nil); err != nil {
		return nil, err
	}

	return mq, nil
}

// SendMessage 发送普通消息
// 注意每次打开一个新通道用完即关闭，如大批量发送消息考虑使用一个已打开的通道
func (mq *MessageQueue) SendMessage(message Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	channel, err := OpenChannel()
	if err != nil {
		log.Printf("打开通道失败：%s\n", err.Error())
		return err
	}
	defer CloseChannel(channel)

	log.Printf("发送消息：%s\n", message.Body)
	err = channel.PublishWithContext(ctx,
		mq.ExchangeName, // exchange
		mq.RouteKey,     // route key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Body),
		},
	)
	return err
}

// SendDelayMessage 发送延迟消息
// 注意每次打开一个新通道用完即关闭，如大批量发送消息考虑使用一个已打开的通道
func (mq *MessageQueue) SendDelayMessage(message Message) error {

	delayQueueName := mq.QueueName + "_delay"
	delayRouteKey := mq.RouteKey + "_delay"

	channel, err := OpenChannel()
	if err != nil {
		log.Printf("打开通道失败：%s\n", err.Error())
		return err
	}
	defer CloseChannel(channel)

	// 定义延迟队列(死信队列)
	dq, err := mq.declareQueue(
		channel,
		delayQueueName,
		amqp.Table{
			"x-dead-letter-exchange":    mq.ExchangeName, // 过期消息exchange
			"x-dead-letter-routing-key": mq.RouteKey,     // 过期消息routing-key
		},
	)
	if err != nil {
		return err
	}

	// 延迟队列绑定到exchange
	if err := mq.bindQueue(channel, dq.Name, delayRouteKey, mq.ExchangeName); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("发送延时消息：%s，%d\n", message.Body, message.DelaySecond)

	// 发送消息，将消息发送到延迟队列，到期后自动路由到正常队列中
	if err := channel.PublishWithContext(ctx,
		mq.ExchangeName,
		delayRouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Body),
			Expiration:  strconv.Itoa(message.DelaySecond * 1000),
		},
	); err != nil {
		return err
	}
	return nil
}

// Consume 获取消费消息
func (mq *MessageQueue) Consume(fn Consumer) {
	channel, err := OpenChannel()
	if err != nil {
		log.Printf("打开通道失败：%s\n", err.Error())
		return
	}
	// defer CloseChannel(channel) // 消费者不能关闭通道

	// 声明队列
	q, err := mq.declareQueue(channel, mq.QueueName, nil)
	if err != nil {
		log.Printf("[队列: %s]声明失败：%s\n", mq.QueueName, err.Error())
		return
	}

	// 队列绑定到exchange
	mq.bindQueue(channel, q.Name, mq.RouteKey, mq.ExchangeName)

	// 设置Qos
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Printf("[队列: %s]Qos设置失败：%s\n", mq.QueueName, err.Error())
		return
	}

	// 监听消息
	msgs, err := channel.Consume(
		q.Name, // queue name,
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("[队列: %s]消费者设置失败：%s\n", mq.QueueName, err.Error())
		return
	}

	go func() {
		for d := range msgs {
			fn(d)
			d.Ack(false)
		}
	}()

	log.Printf("[队列: %s]等待接受消息\n", mq.QueueName)
}

// declareQueue 定义队列
func (mq *MessageQueue) declareQueue(channel *amqp.Channel, name string, args amqp.Table) (amqp.Queue, error) {
	q, err := channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)

	return q, err
}

// declareQueue 定义交换器
func (mq *MessageQueue) declareExchange(channel *amqp.Channel, exchange, kind string, args amqp.Table) error {
	err := channel.ExchangeDeclare(
		exchange,
		kind, // direct, fanout, topic, headers
		true,
		false,
		false,
		false,
		args,
	)
	return err
}

// bindQueue 绑定队列
func (mq *MessageQueue) bindQueue(channel *amqp.Channel, queue, routekey, exchange string) error {
	err := channel.QueueBind(
		queue,
		routekey,
		exchange,
		false,
		nil,
	)
	return err
}
