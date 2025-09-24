package service

import (
	"context"
	"encoding/json"
	"erc20-service/config"
	"erc20-service/internal/db"
	"erc20-service/internal/mq"
	"erc20-service/pkg/logger"
	"fmt"
	"log/slog"
	"math/big"
	"time"
)

// PointsConsumer 积分计算任务消费者
type PointsConsumer struct {
	conn  *mq.Connection
	queue string
	rate  float64
	log   *slog.Logger
}

// NewPointsConsumer 创建消费者
func NewPointsConsumer(conn *mq.Connection, cfg config.RabbitMQConfig, rate float64) *PointsConsumer {
	return &PointsConsumer{
		conn:  conn,
		queue: cfg.Queue,
		rate:  rate,
		log:   logger.New("points-consumer"),
	}
}

// Start 启动消费者
func (c *PointsConsumer) Start(ctx context.Context) error {
	// 注册消费者
	msgs, err := c.conn.Consume(
		c.queue, // 队列名称
		"",      // 消费者标签
		false,   // 自动确认
		false,   // 排他的
		false,   // 不本地
		false,   // 非阻塞
		nil,     // 参数
	)
	if err != nil {
		return fmt.Errorf("注册消费者失败: %v", err)
	}

	c.log.Info("积分计算消费者启动成功")

	// 处理消息
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("消息通道已关闭")
			}

			// 解析任务
			var task mq.PointsCalculationTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				c.log.Warn("解析任务失败", "error", err)
				msg.Nack(false, false)
				continue
			}

			// 计算积分
			if err := c.calculatePoints(task); err != nil {
				c.log.Error("计算积分失败",
					"chain", task.ChainName,
					"user", task.UserAddress,
					"error", err,
				)
				msg.Nack(false, true) // 重新入队
				continue
			}

			// 确认消息
			msg.Ack(false)
		}
	}
}

// 计算用户积分
func (c *PointsConsumer) calculatePoints(task mq.PointsCalculationTask) error {
	c.log.Info("开始计算积分",
		"chain", task.ChainName,
		"user", task.UserAddress,
		"period", fmt.Sprintf("%s - %s", task.PeriodStart, task.PeriodEnd),
	)

	// 1. 获取时间段内的余额变动
	changes, err := db.GetBalanceChangesInPeriod(
		task.ChainName,
		task.UserAddress,
		task.PeriodStart,
		task.PeriodEnd,
	)
	if err != nil {
		return fmt.Errorf("获取余额变动失败: %v", err)
	}

	// 2. 计算积分
	points := c.calculateFromChanges(changes, task.PeriodStart, task.PeriodEnd)
	if points <= 0 {
		c.log.Info("无积分可加", "chain", task.ChainName, "user", task.UserAddress)
		return nil
	}

	// 3. 更新用户总积分
	currentTotal, err := db.GetUserTotalPoints(task.ChainName, task.UserAddress)
	if err != nil {
		return fmt.Errorf("获取当前积分失败: %v", err)
	}

	newTotal := currentTotal + points
	calc := db.PointsCalculation{
		ChainName:    task.ChainName,
		UserAddress:  task.UserAddress,
		PeriodStart:  task.PeriodStart,
		PeriodEnd:    task.PeriodEnd,
		PointsAdded:  points,
		TotalPoints:  newTotal,
		CalculatedAt: time.Now(),
	}

	if err := db.UpdateUserPoints(calc); err != nil {
		return fmt.Errorf("更新积分失败: %v", err)
	}

	c.log.Info("积分计算完成",
		"chain", task.ChainName,
		"user", task.UserAddress,
		"added", points,
		"total", newTotal,
	)
	return nil
}

// 根据余额变动计算积分
func (c *PointsConsumer) calculateFromChanges(changes []db.BalanceChange, start, end time.Time) float64 {
	if len(changes) == 0 {
		return 0
	}

	totalPoints := 0.0
	prevTime := start
	prevBalance := big.NewInt(0) // 初始余额

	totalDuration := end.Sub(start).Hours()
	if totalDuration <= 0 {
		return 0
	}

	for _, change := range changes {
		// 解析当前余额
		currentBalance := new(big.Int)
		currentBalance.SetString(change.BalanceAfter, 10)

		// 计算当前余额的持续时间
		duration := change.EventTime.Sub(prevTime)
		if duration <= 0 {
			prevTime = change.EventTime
			prevBalance = currentBalance
			continue
		}

		// 计算积分：余额 × 0.05 × (持续时间/总周期)
		// 限制余额精度避免溢出
		balanceFloat := new(big.Float).SetInt(prevBalance)
		balanceFloat.Quo(balanceFloat, big.NewFloat(1e18)) // 转换为 ETH 单位
		balance, _ := balanceFloat.Float64()

		periodRatio := duration.Hours() / totalDuration
		points := balance * c.rate * periodRatio
		totalPoints += points

		// 更新状态
		prevTime = change.EventTime
		prevBalance = currentBalance
	}

	// 处理最后一段周期
	lastDuration := end.Sub(prevTime).Hours()
	if lastDuration > 0 {
		balanceFloat := new(big.Float).SetInt(prevBalance)
		balanceFloat.Quo(balanceFloat, big.NewFloat(1e18)) // 转换为 ETH 单位
		balance, _ := balanceFloat.Float64()

		periodRatio := lastDuration / totalDuration
		points := balance * c.rate * periodRatio
		totalPoints += points
	}

	return totalPoints
}
