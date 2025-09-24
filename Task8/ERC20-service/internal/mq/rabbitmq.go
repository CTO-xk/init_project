package mq

import (
	"encoding/json"
	"erc20-service/config"
	"erc20-service/pkg/logger"
	"fmt"
	"log/slog"
	"time"

	"github.com/streadway/amqp"
)

type Connection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	log  *slog.Logger
}

// Consume 包装通道的 Consume，避免暴露内部字段
func (c *Connection) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("连接MQ失败: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建频道失败: %v", err)
	}
	return &Connection{conn: conn, ch: ch, log: logger.New("rabbitmq")}, nil
}
func (c *Connection) Close() error {
	if err := c.ch.Close(); err != nil {
		return fmt.Errorf("关闭频道失败: %v", err)
	}
	return c.conn.Close()
}

// PointsCalculationTask 积分计算任务
type PointsCalculationTask struct {
	ChainName   string    `json:"chain_name"`
	UserAddress string    `json:"user_address"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

// PointsProducer 积分计算任务生产者
type PointsProducer struct {
	ch         *amqp.Channel
	exchange   string
	routingKey string
	log        *slog.Logger
}

// NewPointsProducer 创建积分任务生产者
func NewPointsProducer(conn *Connection, cfg config.RabbitMQConfig) *PointsProducer {
	//声明交换机
	err := conn.ch.ExchangeDeclare(
		cfg.Exchange, // 交换机名称
		"direct",     // 类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部的
		false,        // 非阻塞
		nil,          // 参数
	)
	if err != nil {
		panic(fmt.Sprintf("声明交换机失败: %v", err))
	}
	//声明队列
	_, err = conn.ch.QueueDeclare(
		cfg.Queue, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他的
		false,     // 非阻塞
		nil,       // 参数
	)
	if err != nil {
		panic(fmt.Sprintf("声明队列失败: %v", err))
	}
	//绑定队列到交换机
	err = conn.ch.QueueBind(
		cfg.Queue,      // 队列名称
		cfg.RoutingKey, // 路由键
		cfg.Exchange,   // 交换机名称
		false,          // 是否持久化
		nil,            // 参数
	)
	if err != nil {
		panic(fmt.Sprintf("绑定队列到交换机失败: %v", err))
	}
	return &PointsProducer{
		ch:         conn.ch,
		exchange:   cfg.Exchange,
		routingKey: cfg.RoutingKey,
		log:        logger.New("points-producer"),
	}
}

// Publish 发布积分计算任务
func (p *PointsProducer) Publish(task PointsCalculationTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %v", err)
	}
	return p.ch.Publish(
		p.exchange,   // 交换机
		p.routingKey, // 路由键
		false,        // 强制的
		false,        // 立即的
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)
	if err != nil {
		return fmt.Errorf("发布任务失败: %v", err)
	}
	p.log.Debug("发布积分计算任务",
		"chain", task.ChainName,
		"user", task.UserAddress,
	)
	return nil
}
