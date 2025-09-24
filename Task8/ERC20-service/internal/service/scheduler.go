package service

import (
	"context"
	"erc20-service/internal/db"
	"erc20-service/internal/mq"
	"erc20-service/pkg/logger"
	"fmt"
	"log/slog"
	"time"
)

// Scheduler 积分计算定时调度器
type Scheduler struct {
	interval int // 调度间隔（分钟）
	chains   []string
	producer *mq.PointsProducer
	log      *slog.Logger
}

// NewScheduler 创建调度器
func NewScheduler(interval int, chains []string, producer *mq.PointsProducer) *Scheduler {
	return &Scheduler{
		interval: interval,
		chains:   chains,
		producer: producer,
		log:      logger.New("scheduler"),
	}
}

// Start 启动调度器
func (s *Scheduler) Start(ctx context.Context) {
	s.log.Info("启动积分计算调度器", "interval", s.interval, "unit", "minutes")

	// 立即执行一次
	s.scheduleAllChains()

	// 定时执行
	ticker := time.NewTicker(time.Duration(s.interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("调度器已停止")
			return
		case <-ticker.C:
			s.scheduleAllChains()
		}
	}
}

// 为所有链调度积分计算任务
func (s *Scheduler) scheduleAllChains() {
	s.log.Info("开始调度积分计算任务")

	for _, chain := range s.chains {
		if err := s.scheduleChain(chain); err != nil {
			s.log.Error("调度链积分任务失败", "chain", chain, "error", err)
		}
	}
}

// 为单个链调度积分计算任务
func (s *Scheduler) scheduleChain(chainName string) error {
	// 获取链上所有用户
	users, err := db.GetUsersByChain(chainName)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %v", err)
	}

	s.log.Info("调度积分计算任务", "chain", chainName, "user_count", len(users))

	// 计算周期
	now := time.Now()
	periodEnd := now
	periodStart := now.Add(-time.Duration(s.interval) * time.Minute)

	// 为每个用户创建任务
	for _, user := range users {
		// 获取用户上次计算时间
		lastCalc, err := db.GetUserLastCalculatedTime(chainName, user)
		if err != nil {
			s.log.Warn("获取用户上次计算时间失败", "user", user, "error", err)
			continue
		}

		// 检查是否需要回溯计算
		hoursSinceLastCalc := now.Sub(lastCalc).Hours()

		// 智能检测积分缺失：
		// 1. 如果上次计算时间早于当前周期开始时间，说明有积分缺失
		// 2. 如果超过正常计算间隔的2倍时间未计算，也认为有缺失
		normalInterval := time.Duration(s.interval) * time.Minute
		expectedLastCalc := now.Add(-normalInterval)

		needsBackfill := lastCalc.Before(periodStart) || lastCalc.Before(expectedLastCalc)

		if needsBackfill {
			s.log.Warn("检测到积分计算缺失，启动回溯",
				"user", user,
				"last_calc", lastCalc,
				"period_start", periodStart,
				"expected_last_calc", expectedLastCalc,
				"hours_behind", hoursSinceLastCalc)

			// 执行回溯计算（从上次计算时间到当前周期开始）
			backfillEnd := periodStart
			if err := s.backfillUserPoints(chainName, user, lastCalc, backfillEnd); err != nil {
				s.log.Error("回溯计算失败", "user", user, "error", err)
				// 即使回溯失败，也继续正常计算当前周期
			}
		}

		// 确定实际计算周期
		actualStart := lastCalc
		if actualStart.Before(periodStart) {
			actualStart = periodStart
		}

		// 发布任务
		task := mq.PointsCalculationTask{
			ChainName:   chainName,
			UserAddress: user,
			PeriodStart: actualStart,
			PeriodEnd:   periodEnd,
		}

		if err := s.producer.Publish(task); err != nil {
			s.log.Warn("发布任务失败", "user", user, "error", err)
		}
	}

	return nil
}

// 回溯计算用户积分
func (s *Scheduler) backfillUserPoints(chainName, userAddr string, startTime, endTime time.Time) error {
	s.log.Info("开始回溯计算用户积分", "chain", chainName, "user", userAddr, "start", startTime, "end", endTime)

	// 根据时间跨度选择合适的分割粒度
	duration := endTime.Sub(startTime)
	var splitDuration time.Duration

	if duration <= time.Hour {
		// 小于1小时，按15分钟分割
		splitDuration = 15 * time.Minute
	} else if duration <= 24*time.Hour {
		// 小于24小时，按1小时分割
		splitDuration = time.Hour
	} else {
		// 大于24小时，按6小时分割
		splitDuration = 6 * time.Hour
	}

	periods := s.splitTimeRange(startTime, endTime, splitDuration)
	taskCount := 0

	for _, period := range periods {
		// 检查该时间段是否已经计算过
		hasCalculated, err := db.HasPointsCalculated(chainName, userAddr, period.Start, period.End)
		if err != nil {
			s.log.Warn("检查计算状态失败", "user", userAddr, "period", period, "error", err)
			continue
		}

		if hasCalculated {
			s.log.Debug("时间段已计算", "user", userAddr, "period", period)
			continue
		}

		// 发布计算任务
		task := mq.PointsCalculationTask{
			ChainName:   chainName,
			UserAddress: userAddr,
			PeriodStart: period.Start,
			PeriodEnd:   period.End,
		}

		if err := s.producer.Publish(task); err != nil {
			s.log.Warn("发布回溯任务失败", "user", userAddr, "period", period, "error", err)
			continue
		}

		taskCount++
	}

	s.log.Info("回溯任务发布完成", "user", userAddr, "task_count", taskCount, "duration", duration)
	return nil
}

// 分割时间范围
func (s *Scheduler) splitTimeRange(start, end time.Time, duration time.Duration) []db.TimePeriod {
	var periods []db.TimePeriod

	current := start
	for current.Before(end) {
		periodEnd := current.Add(duration)
		if periodEnd.After(end) {
			periodEnd = end
		}

		periods = append(periods, db.TimePeriod{
			Start: current,
			End:   periodEnd,
		})

		current = periodEnd
	}

	return periods
}
