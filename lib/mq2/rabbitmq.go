/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-09 15:19:31
 * @LastEditTime: 2024-07-31 10:15:21
 */
package mq2

import (
	"context"
	"fmt"
	"iris-project/app/config"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 消息体：DelayTime 仅在 SendDelayMessage 方法有效
type Message struct {
	DelayTime int // desc:延迟时间(秒)
	Body      string
}

type MessageQueue struct {
	conn         *amqp.Connection // amqp链接对象
	ch           *amqp.Channel    // channel对象
	ExchangeName string           // 交换器名称
	ExchangeKind string           // 交换器类型 direct, fanout, topic, headers
	RouteKey     string           // 路由名称
	QueueName    string           // 队列名称
}

// 消费者回调方法
type Consumer func(amqp.Delivery)

// NewRabbitMQ 新建 rabbitmq 实例，建立连接声明交换机
func NewRabbitMQ(exchangeName, exchangeKind, routeKey, queueName string) *MessageQueue {
	var messageQueue = &MessageQueue{
		ExchangeName: exchangeName,
		ExchangeKind: exchangeKind,
		RouteKey:     routeKey,
		QueueName:    queueName,
	}

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.RabbitMQ.Username,
		config.RabbitMQ.Password,
		config.RabbitMQ.IP,
		config.RabbitMQ.Port,
	)

	// 建立amqp链接
	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	messageQueue.conn = conn

	// 建立channel通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	messageQueue.ch = ch

	// 声明exchange交换器
	messageQueue.declareExchange(exchangeName, exchangeKind, nil)

	return messageQueue
}

// SendMessage 发送普通消息
func (mq *MessageQueue) SendMessage(message Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mq.ch.PublishWithContext(ctx,
		mq.ExchangeName, // exchange
		mq.RouteKey,     // route key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Body),
		},
	)
	failOnError(err, "send common msg err")
}

// SendDelayMessage 发送延迟消息
func (mq *MessageQueue) SendDelayMessage(message Message) {
	// delayQueueName := mq.QueueName + "_delay:" + strconv.Itoa(message.DelayTime)
	// delayRouteKey := mq.RouteKey + "_delay:" + strconv.Itoa(message.DelayTime)

	delayQueueName := mq.QueueName + "_delay"
	delayRouteKey := mq.RouteKey + "_delay"

	// 定义延迟队列(死信队列)
	dq := mq.declareQueue(
		delayQueueName,
		amqp.Table{
			"x-dead-letter-exchange":    mq.ExchangeName, // 指定死信交换机
			"x-dead-letter-routing-key": mq.RouteKey,     // 指定死信routing-key
		},
	)

	// 延迟队列绑定到exchange
	mq.bindQueue(dq.Name, delayRouteKey, mq.ExchangeName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息，将消息发送到延迟队列，到期后自动路由到正常队列中
	err := mq.ch.PublishWithContext(ctx,
		mq.ExchangeName,
		delayRouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Body),
			Expiration:  strconv.Itoa(message.DelayTime * 1000),
		},
	)
	failOnError(err, "send delay msg err")
}

// Consume 获取消费消息
func (mq *MessageQueue) Consume(fn Consumer) {
	// 声明队列
	q := mq.declareQueue(mq.QueueName, nil)

	// 队列绑定到exchange
	mq.bindQueue(q.Name, mq.RouteKey, mq.ExchangeName)

	// 设置Qos
	err := mq.ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	// 监听消息
	msgs, err := mq.ch.Consume(
		q.Name, // queue name,
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// forever := make(chan bool), 注册在主进程，不需要阻塞

	go func() {
		for d := range msgs {
			fn(d)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	// <-forever
}

// Close 关闭链接
func (mq *MessageQueue) Close() {
	mq.ch.Close()
	mq.conn.Close()
}

// declareQueue 定义队列
func (mq *MessageQueue) declareQueue(name string, args amqp.Table) amqp.Queue {
	q, err := mq.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare a delay_queue")

	return q
}

// declareQueue 定义交换器
func (mq *MessageQueue) declareExchange(exchange, kind string, args amqp.Table) {
	err := mq.ch.ExchangeDeclare(
		exchange,
		kind, // direct, fanout, topic, headers
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare an exchange")
}

// bindQueue 绑定队列
func (mq *MessageQueue) bindQueue(queue, routekey, exchange string) {
	err := mq.ch.QueueBind(
		queue,
		routekey,
		exchange,
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
}

// failOnError 错误处理
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}
