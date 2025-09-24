package backfill

import (
	"erc20-service/cmd"
	"erc20-service/config"
	"erc20-service/internal/db"
	"erc20-service/internal/mq"
	"erc20-service/pkg/logger"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	backfillCmd = &cobra.Command{
		Use:   "backfill",
		Short: "积分回溯计算工具",
		Long:  "用于处理服务中断期间的积分回溯计算",
	}

	backfillPointsCmd = &cobra.Command{
		Use:   "points [chain_name] [start_time] [end_time]",
		Short: "回溯计算指定时间段的积分",
		Long: `回溯计算指定时间段的积分
参数:
  chain_name: 链名称 (如: sepolia)
  start_time: 开始时间 (格式: 2006-01-02T15:04:05Z)
  end_time: 结束时间 (格式: 2006-01-02T15:04:05Z)
  
示例:
  ./erc20-service backfill points sepolia 2024-01-01T00:00:00Z 2024-01-03T00:00:00Z`,
		Args: cobra.ExactArgs(3),
		Run:  runBackfillPoints,
	}

	backfillCheckCmd = &cobra.Command{
		Use:   "check [chain_name]",
		Short: "检查积分计算状态",
		Long:  "检查指定链的积分计算状态，显示缺失的时间段",
		Args:  cobra.ExactArgs(1),
		Run:   runBackfillCheck,
	}

	backfillScanCmd = &cobra.Command{
		Use:   "scan [chain_name]",
		Short: "扫描并修复所有积分缺失",
		Long:  "扫描指定链的所有用户，自动修复积分计算缺失",
		Args:  cobra.ExactArgs(1),
		Run:   runBackfillScan,
	}

	log = logger.New("backfill")
)

func init() {
	cmd.RootCmd.AddCommand(backfillCmd)
	backfillCmd.AddCommand(backfillPointsCmd)
	backfillCmd.AddCommand(backfillCheckCmd)
	backfillCmd.AddCommand(backfillScanCmd)
}

func runBackfillPoints(cmd *cobra.Command, args []string) {
	chainName := args[0]
	startTimeStr := args[1]
	endTimeStr := args[2]

	// 解析时间
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		logger.Fatal("解析开始时间失败", "error", err)
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		logger.Fatal("解析结束时间失败", "error", err)
	}

	if startTime.After(endTime) {
		logger.Fatal("开始时间不能晚于结束时间")
	}

	// 加载配置
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("加载配置失败", "error", err)
	}

	// 初始化数据库
	if err := db.Init(cfg.Database); err != nil {
		logger.Fatal("初始化数据库失败", "error", err)
	}

	// 初始化MQ连接
	mqConn, err := mq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Fatal("初始化MQ连接失败", "error", err)
	}
	defer mqConn.Close()

	// 创建积分生产者
	pointsProducer := mq.NewPointsProducer(mqConn, cfg.RabbitMQ)

	// 执行回溯计算
	if err := backfillPointsForChain(chainName, startTime, endTime, pointsProducer, cfg.Points.Rate); err != nil {
		logger.Fatal("回溯计算失败", "error", err)
	}

	log.Info("回溯计算完成", "chain", chainName, "start", startTime, "end", endTime)
}

func runBackfillCheck(cmd *cobra.Command, args []string) {
	chainName := args[0]

	// 加载配置
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("加载配置失败", "error", err)
	}

	// 初始化数据库
	if err := db.Init(cfg.Database); err != nil {
		logger.Fatal("初始化数据库失败", "error", err)
	}

	// 检查积分计算状态
	if err := checkPointsCalculationStatus(chainName); err != nil {
		logger.Fatal("检查失败", "error", err)
	}
}

// 回溯计算指定链的积分
func backfillPointsForChain(chainName string, startTime, endTime time.Time, producer *mq.PointsProducer, rate float64) error {
	log.Info("开始回溯计算", "chain", chainName, "start", startTime, "end", endTime)

	// 获取链上所有用户
	users, err := db.GetUsersByChain(chainName)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %v", err)
	}

	if len(users) == 0 {
		log.Info("链上无用户", "chain", chainName)
		return nil
	}

	log.Info("找到用户", "count", len(users))

	// 按小时分割时间段，避免单次计算时间过长
	periods := splitTimeRange(startTime, endTime, time.Hour)

	totalTasks := 0
	for _, period := range periods {
		for _, user := range users {
			// 检查该时间段是否已经计算过
			hasCalculated, err := db.HasPointsCalculated(chainName, user, period.Start, period.End)
			if err != nil {
				log.Warn("检查计算状态失败", "user", user, "period", period, "error", err)
				continue
			}

			if hasCalculated {
				log.Debug("时间段已计算", "user", user, "period", period)
				continue
			}

			// 发布计算任务
			task := mq.PointsCalculationTask{
				ChainName:   chainName,
				UserAddress: user,
				PeriodStart: period.Start,
				PeriodEnd:   period.End,
			}

			if err := producer.Publish(task); err != nil {
				log.Warn("发布任务失败", "user", user, "period", period, "error", err)
				continue
			}

			totalTasks++
		}
	}

	log.Info("回溯任务发布完成", "total_tasks", totalTasks)
	return nil
}

// 检查积分计算状态
func checkPointsCalculationStatus(chainName string) error {
	log.Info("检查积分计算状态", "chain", chainName)

	// 获取链上所有用户
	users, err := db.GetUsersByChain(chainName)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %v", err)
	}

	if len(users) == 0 {
		log.Info("链上无用户", "chain", chainName)
		return nil
	}

	// 检查每个用户的积分计算状态
	for _, user := range users {
		lastCalc, err := db.GetUserLastCalculatedTime(chainName, user)
		if err != nil {
			log.Warn("获取用户计算时间失败", "user", user, "error", err)
			continue
		}

		now := time.Now()
		hoursSinceLastCalc := now.Sub(lastCalc).Hours()

		if hoursSinceLastCalc > 2 { // 超过2小时未计算
			log.Warn("用户积分计算滞后",
				"user", user,
				"last_calc", lastCalc,
				"hours_behind", hoursSinceLastCalc)
		} else {
			log.Info("用户积分计算正常", "user", user, "last_calc", lastCalc)
		}
	}

	return nil
}

func runBackfillScan(cmd *cobra.Command, args []string) {
	chainName := args[0]

	// 加载配置
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("加载配置失败", "error", err)
	}

	// 初始化数据库
	if err := db.Init(cfg.Database); err != nil {
		logger.Fatal("初始化数据库失败", "error", err)
	}

	// 初始化MQ连接
	mqConn, err := mq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Fatal("初始化MQ连接失败", "error", err)
	}
	defer mqConn.Close()

	// 创建积分生产者
	pointsProducer := mq.NewPointsProducer(mqConn, cfg.RabbitMQ)

	// 执行扫描和修复
	if err := scanAndFixMissingPoints(chainName, pointsProducer, cfg.Points.Rate, cfg.Points.Interval); err != nil {
		logger.Fatal("扫描修复失败", "error", err)
	}

	log.Info("扫描修复完成", "chain", chainName)
}

// 扫描并修复积分缺失
func scanAndFixMissingPoints(chainName string, producer *mq.PointsProducer, rate float64, intervalMinutes int) error {
	log.Info("开始扫描积分缺失", "chain", chainName)

	// 获取链上所有用户
	users, err := db.GetUsersByChain(chainName)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %v", err)
	}

	if len(users) == 0 {
		log.Info("链上无用户", "chain", chainName)
		return nil
	}

	log.Info("找到用户", "count", len(users))

	now := time.Now()
	normalInterval := time.Duration(intervalMinutes) * time.Minute
	fixedCount := 0
	totalTasks := 0

	for _, user := range users {
		// 获取用户上次计算时间
		lastCalc, err := db.GetUserLastCalculatedTime(chainName, user)
		if err != nil {
			log.Warn("获取用户计算时间失败", "user", user, "error", err)
			continue
		}

		// 检查是否需要修复
		expectedLastCalc := now.Add(-normalInterval)
		needsFix := lastCalc.Before(expectedLastCalc)

		if needsFix {
			log.Info("发现用户积分缺失",
				"user", user,
				"last_calc", lastCalc,
				"expected", expectedLastCalc)

			// 计算修复的时间范围
			backfillEnd := now.Add(-normalInterval) // 修复到正常周期开始

			// 按小时分割时间段
			periods := splitTimeRange(lastCalc, backfillEnd, time.Hour)

			userTasks := 0
			for _, period := range periods {
				// 检查该时间段是否已经计算过
				hasCalculated, err := db.HasPointsCalculated(chainName, user, period.Start, period.End)
				if err != nil {
					log.Warn("检查计算状态失败", "user", user, "period", period, "error", err)
					continue
				}

				if hasCalculated {
					continue
				}

				// 发布计算任务
				task := mq.PointsCalculationTask{
					ChainName:   chainName,
					UserAddress: user,
					PeriodStart: period.Start,
					PeriodEnd:   period.End,
				}

				if err := producer.Publish(task); err != nil {
					log.Warn("发布修复任务失败", "user", user, "period", period, "error", err)
					continue
				}

				userTasks++
				totalTasks++
			}

			if userTasks > 0 {
				fixedCount++
				log.Info("用户修复任务发布完成", "user", user, "task_count", userTasks)
			}
		}
	}

	log.Info("扫描修复完成",
		"total_users", len(users),
		"fixed_users", fixedCount,
		"total_tasks", totalTasks)

	return nil
}

// 时间段结构
type TimePeriod struct {
	Start time.Time
	End   time.Time
}

// 分割时间范围
func splitTimeRange(start, end time.Time, duration time.Duration) []TimePeriod {
	var periods []TimePeriod

	current := start
	for current.Before(end) {
		periodEnd := current.Add(duration)
		if periodEnd.After(end) {
			periodEnd = end
		}

		periods = append(periods, TimePeriod{
			Start: current,
			End:   periodEnd,
		})

		current = periodEnd
	}

	return periods
}
