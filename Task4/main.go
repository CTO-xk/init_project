package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"sepolia-ops/block"
	"sepolia-ops/client"
	"sepolia-ops/config"
	"sepolia-ops/logger"
	"sepolia-ops/transaction"
)

func main() {
	// 1. 解析命令行参数
	configPath := flag.String("config", "", "配置文件路径")
	flag.Parse()

	// 2. 加载配置
	cfg := loadConfig(*configPath)

	// 验证配置
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("配置验证失败: %v", err))
	}

	// 3. 初始化日志系统
	if err := logger.InitLogger(cfg.Env, cfg.LogPath); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync() // 确保程序退出前刷新日志

	logger.L().Info("应用启动",
		zap.String("env", cfg.Env),
		zap.String("log_path", cfg.LogPath))

	// 4. 初始化以太坊客户端
	ethClient, err := client.NewETHClient(cfg.InfuraAPIKey)
	if err != nil {
		logger.L().Fatal("客户端初始化失败", zap.Error(err))
	}
	defer ethClient.Close()
	logger.L().Info("已连接到Sepolia测试网")

	// 5. 区块查询流程
	blockService := block.NewBlockService(ethClient.GetClient())
	ctx := context.Background()

	// 查询指定区块
	logger.L().Info("开始区块查询", zap.Uint64("block_number", cfg.TargetBlock.Uint64()))
	targetBlock, err := blockService.GetBlockByNumber(ctx, cfg.TargetBlock)
	if err != nil {
		logger.L().Fatal("区块查询失败", zap.Error(err))
	}
	blockService.PrintBlockInfo(targetBlock)

	// 6. 交易发送流程
	chainID, err := ethClient.GetChainID(ctx)
	if err != nil {
		logger.L().Fatal("获取链ID失败", zap.Error(err))
	}

	txService, err := transaction.NewTransactionService(
		ethClient.GetClient(),
		cfg.PrivateKey,
		chainID,
	)
	if err != nil {
		logger.L().Fatal("交易服务初始化失败", zap.Error(err))
	}

	logger.L().Info("开始创建并发送交易",
		zap.String("recipient", cfg.RecipientAddr.Hex()),
		zap.String("amount_wei", cfg.TransferWei.String()))

	signedTx, err := txService.CreateAndSendTransaction(
		ctx,
		cfg.RecipientAddr,
		cfg.TransferWei,
	)
	if err != nil {
		logger.L().Fatal("交易处理失败", zap.Error(err))
	}
	txService.PrintTransactionInfo(signedTx)

	// 7. 优雅退出
	waitForShutdown()
}

// loadConfig 加载配置（支持从文件或默认值）
func loadConfig(configPath string) config.Config {
	// 实际项目中可以从JSON/YAML文件加载
	cfg := config.DefaultConfig()

	// 示例：手动设置配置（实际项目建议从配置文件读取）
	cfg.InfuraAPIKey = "YOUR_INFURA_API_KEY"                      // 替换为你的Infura API Key
	cfg.PrivateKey = "YOUR_PRIVATE_KEY"                           // 替换为发送方私钥（不含0x）
	cfg.RecipientAddr = common.HexToAddress("0xRecipientAddress") // 替换为接收方地址
	cfg.TransferWei = big.NewInt(100000000000000000)              // 转账0.1 ETH（1e17 Wei）
	cfg.TargetBlock = big.NewInt(5000000)                         // 目标查询区块号

	return cfg
}

// waitForShutdown 等待系统信号优雅退出
func waitForShutdown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	logger.L().Info("接收到退出信号，正在优雅退出...")
}
