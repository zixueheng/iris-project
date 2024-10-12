/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2024-07-30 10:48:59
 * @LastEditTime: 2024-08-02 11:14:14
 */
package middleware

import (
	"iris-project/lib/mq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MQ_PREFIX = "iris."

	DIRECT_EXCHANGE  = MQ_PREFIX + "direct.exchange" // 直连交换机
	DIRECT_QUEUE     = MQ_PREFIX + "direct.queue"    // 直连队列
	DIRECT_ROUTE_KEY = MQ_PREFIX + "direct.route.key"

	// FANOUT_EXCHANGE  = MQ_PREFIX + "fanout.exchange" // 扇型交换机
	// FANOUT_QUEUE     = MQ_PREFIX + "fanout.queue"    // 扇型队列
	// FANOUT_ROUTE_KEY = MQ_PREFIX + "fanout.route.key"

	// TOPIC_EXCHANGE      = MQ_PREFIX + "topic.exchange"  // 主题交换器
	// TOPIC_QUEUE_ONE     = MQ_PREFIX + "topic.queue.one" // 主题队列1
	// TOPIC_ROUTE_KEY_ONE = MQ_PREFIX + "topic.route.key.one"
	// TOPIC_QUEUE_TWO     = MQ_PREFIX + "topic.queue.two" // 主题队列2
	// TOPIC_ROUTE_KEY_TWO = MQ_PREFIX + "topic.route.key.two"
)

var (
	MqDirect *mq.MessageQueue

	// mqFanout *mq.MessageQueue
	// mqTopic  *mq.MessageQueue
)

func InitMq() {
	var err error
	MqDirect, err = mq.New(DIRECT_EXCHANGE, mq.DirectType, DIRECT_ROUTE_KEY, DIRECT_QUEUE, true)
	if err != nil {
		log.Fatalf("rabbit mq init error: %s", err.Error())
		return
	}
	for i := 1; i <= 3; i++ { // 启动三个消费者
		go func(i int) {
			// 执行消费
			MqDirect.Consume(func(d amqp.Delivery) {
				log.Printf("消费者%d收到消息：%s\n", i, string(d.Body))
			})
		}(i)
	}
}
