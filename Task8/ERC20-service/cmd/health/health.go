package health

import (
	"erc20-service/cmd"
	"erc20-service/config"
	"erc20-service/internal/db"
	"erc20-service/pkg/logger"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "系统健康检查",
		Long:  "检查系统各组件的健康状态",
	}

	healthCheckCmd = &cobra.Command{
		Use:   "check",
		Short: "执行健康检查",
		Run:   runHealthCheck,
	}

	log = logger.New("health")
)

func init() {
	cmd.RootCmd.AddCommand(healthCmd)
	healthCmd.AddCommand(healthCheckCmd)
}

func runHealthCheck(cmd *cobra.Command, args []string) {
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

	// 执行健康检查
	status := performHealthCheck(cfg)

	// 输出结果
	if status.IsHealthy {
		log.Info("系统健康检查通过", "details", status)
	} else {
		log.Error("系统健康检查失败", "details", status)
	}
}

// HealthStatus 健康状态
type HealthStatus struct {
	IsHealthy               bool                    `json:"is_healthy"`
	DatabaseStatus          ComponentStatus         `json:"database_status"`
	ChainStatuses           map[string]ChainStatus  `json:"chain_statuses"`
	PointsCalculationStatus PointsCalculationStatus `json:"points_calculation_status"`
	OverallScore            int                     `json:"overall_score"`
}

// ComponentStatus 组件状态
type ComponentStatus struct {
	IsHealthy bool   `json:"is_healthy"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
}

// ChainStatus 链状态
type ChainStatus struct {
	IsHealthy          bool      `json:"is_healthy"`
	LastProcessedBlock uint64    `json:"last_processed_block"`
	LastProcessedTime  time.Time `json:"last_processed_time"`
	HoursBehind        float64   `json:"hours_behind"`
	Message            string    `json:"message"`
}

// PointsCalculationStatus 积分计算状态
type PointsCalculationStatus struct {
	IsHealthy          bool    `json:"is_healthy"`
	TotalUsers         int     `json:"total_users"`
	UsersBehind        int     `json:"users_behind"`
	AverageHoursBehind float64 `json:"average_hours_behind"`
	Message            string  `json:"message"`
}

// 执行健康检查
func performHealthCheck(cfg *config.Config) HealthStatus {
	status := HealthStatus{
		IsHealthy:     true,
		ChainStatuses: make(map[string]ChainStatus),
	}

	// 检查数据库
	status.DatabaseStatus = checkDatabase()
	if !status.DatabaseStatus.IsHealthy {
		status.IsHealthy = false
	}

	// 检查各链状态
	for _, chain := range cfg.Chains {
		chainStatus := checkChainStatus(chain.Name)
		status.ChainStatuses[chain.Name] = chainStatus
		if !chainStatus.IsHealthy {
			status.IsHealthy = false
		}
	}

	// 检查积分计算状态
	status.PointsCalculationStatus = checkPointsCalculationStatus()
	if !status.PointsCalculationStatus.IsHealthy {
		status.IsHealthy = false
	}

	// 计算总体评分
	status.OverallScore = calculateOverallScore(status)

	return status
}

// 检查数据库状态
func checkDatabase() ComponentStatus {
	// 测试数据库连接
	if err := db.DB.Ping(); err != nil {
		return ComponentStatus{
			IsHealthy: false,
			Message:   "数据库连接失败",
			Details:   err.Error(),
		}
	}

	// 检查关键表是否存在
	tables := []string{"chain_status", "user_balances", "balance_changes", "user_points", "points_calculation_history"}
	for _, table := range tables {
		var count int
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s LIMIT 1", table)).Scan(&count)
		if err != nil {
			return ComponentStatus{
				IsHealthy: false,
				Message:   "数据库表检查失败",
				Details:   fmt.Sprintf("表 %s 不可访问: %v", table, err),
			}
		}
	}

	return ComponentStatus{
		IsHealthy: true,
		Message:   "数据库连接正常",
	}
}

// 检查链状态
func checkChainStatus(chainName string) ChainStatus {
	// 获取最后处理的区块
	lastBlock, err := db.GetLastProcessedBlock(chainName)
	if err != nil {
		return ChainStatus{
			IsHealthy: false,
			Message:   "获取链状态失败",
		}
	}

	// 获取最后处理时间（从chain_status表）
	var lastProcessedTime time.Time
	err = db.QueryRow(`
		SELECT updated_at FROM chain_status WHERE chain_name = ?
	`, chainName).Scan(&lastProcessedTime)
	if err != nil {
		return ChainStatus{
			IsHealthy: false,
			Message:   "获取最后处理时间失败",
		}
	}

	// 计算滞后时间
	hoursBehind := time.Since(lastProcessedTime).Hours()

	isHealthy := hoursBehind < 2 // 超过2小时认为不健康

	return ChainStatus{
		IsHealthy:          isHealthy,
		LastProcessedBlock: lastBlock,
		LastProcessedTime:  lastProcessedTime,
		HoursBehind:        hoursBehind,
		Message:            fmt.Sprintf("最后处理区块: %d, 滞后: %.2f小时", lastBlock, hoursBehind),
	}
}

// 检查积分计算状态
func checkPointsCalculationStatus() PointsCalculationStatus {
	// 获取所有用户
	users, err := db.GetUsersByChain("sepolia") // 假设检查sepolia链
	if err != nil {
		return PointsCalculationStatus{
			IsHealthy: false,
			Message:   "获取用户列表失败",
		}
	}

	totalUsers := len(users)
	usersBehind := 0
	totalHoursBehind := 0.0

	now := time.Now()
	for _, user := range users {
		lastCalc, err := db.GetUserLastCalculatedTime("sepolia", user)
		if err != nil {
			continue
		}

		hoursBehind := now.Sub(lastCalc).Hours()
		if hoursBehind > 2 { // 超过2小时认为滞后
			usersBehind++
			totalHoursBehind += hoursBehind
		}
	}

	avgHoursBehind := 0.0
	if usersBehind > 0 {
		avgHoursBehind = totalHoursBehind / float64(usersBehind)
	}

	isHealthy := usersBehind == 0 || avgHoursBehind < 24 // 平均滞后小于24小时认为健康

	return PointsCalculationStatus{
		IsHealthy:          isHealthy,
		TotalUsers:         totalUsers,
		UsersBehind:        usersBehind,
		AverageHoursBehind: avgHoursBehind,
		Message:            fmt.Sprintf("总用户: %d, 滞后用户: %d, 平均滞后: %.2f小时", totalUsers, usersBehind, avgHoursBehind),
	}
}

// 计算总体评分
func calculateOverallScore(status HealthStatus) int {
	score := 100

	// 数据库状态
	if !status.DatabaseStatus.IsHealthy {
		score -= 50
	}

	// 链状态
	for _, chainStatus := range status.ChainStatuses {
		if !chainStatus.IsHealthy {
			score -= 20
		}
		if chainStatus.HoursBehind > 24 {
			score -= 10
		}
	}

	// 积分计算状态
	if !status.PointsCalculationStatus.IsHealthy {
		score -= 30
	}
	if status.PointsCalculationStatus.AverageHoursBehind > 48 {
		score -= 20
	}

	if score < 0 {
		score = 0
	}

	return score
}
