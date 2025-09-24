package daemon

import (
	"context"
	rootcmd "erc20-service/cmd"
	"erc20-service/config"
	"erc20-service/internal/chain"
	"erc20-service/internal/db"
	"erc20-service/internal/mq"
	"erc20-service/internal/service"
	"erc20-service/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "启动ERC20事件追踪守护进程",
		Long:  "启动多链事件监听、余额更新、积分计算调度等后台任务",
		Run:   runDaemon,
	}
	log = logger.New("daemon")
)

func init() {
	rootcmd.RootCmd.AddCommand(daemonCmd)
}

func runDaemon(cmd *cobra.Command, args []string) {
	// 1. 加载配置
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("加载配置失败", "error", err)
	}

	// 2. 初始化基础设施
	// 2.1 数据库连接
	if err := db.Init(cfg.Database); err != nil {
		logger.Fatal("初始化数据库失败", "error", err)
	}
	// 2.2 初始化链状态
	if err := db.InitChainStatus(cfg.Chains); err != nil {
		logger.Fatal("初始化链状态失败", "error", err)
	}
	// 2.3 载入内置 ERC20 ABI
	abiBytes := chain.ERC20ABI
	// 2.4 初始化MQ连接
	mqConn, err := mq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Fatal("初始化MQ连接失败", "error", err)
	}
	defer mqConn.Close()

	// 3. 创建核心服务组件
	// 3.1 MQ生产者（发送积分计算任务）
	pointsProducer := mq.NewPointsProducer(mqConn, cfg.RabbitMQ)
	// 3.2 MQ消费者（处理积分计算）
	pointsConsumer := service.NewPointsConsumer(mqConn, cfg.RabbitMQ, cfg.Points.Rate)
	// 3.3 多链事件监听管理器
	chainManager := chain.NewManager(cfg.Chains, abiBytes, pointsProducer)
	// 3.4 积分计算定时调度器
	scheduler := service.NewScheduler(cfg.Points.Interval, chainManager.GetChainNames(), pointsProducer)

	// 4. 启动服务组件
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// 4.1 启动事件监听器
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := chainManager.Start(ctx); err != nil {
			log.Error("链管理器退出", "error", err)
		}
	}()

	// 4.2 启动积分计算消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := pointsConsumer.Start(ctx); err != nil {
			log.Error("积分消费者退出", "error", err)
		}
	}()

	// 4.3 启动定时调度器
	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start(ctx)
		log.Info("积分调度器退出")
	}()

	// 5. 等待退出信号
	log.Info("服务启动成功，等待退出信号...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigChan:
		log.Info("收到退出信号", "signal", sig.String())
	case err := <-chainManager.ExitChan():
		log.Error("监听器异常退出", "error", err)
	}

	// 6. 优雅关闭
	log.Info("开始优雅关闭服务...")
	cancel()
	wg.Wait()
	log.Info("所有服务已关闭，退出程序")
}
