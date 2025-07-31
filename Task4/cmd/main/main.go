package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/CTO-xk/init_project/Task4/pkg/block"
	"github.com/CTO-xk/init_project/Task4/pkg/client"
	"github.com/CTO-xk/init_project/Task4/pkg/config"
	"github.com/CTO-xk/init_project/Task4/pkg/logger"
	"github.com/CTO-xk/init_project/Task4/pkg/transaction"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	// 解析配置文件路径参数
	configPath := flag.String("config", "", "pkg/config/config.yaml）")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化日志系统
	if err := logger.Init(cfg.Log); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync() // 程序退出时同步日志

	// 创建以太坊客户端连接
	// 使用不需要认证的公共Sepolia RPC节点
	rpcURL := "https://eth-sepolia.g.alchemy.com/v2/XE2nmOMCIb6XVkP4Rj7Ar"
	ethClient, err := client.NewEthClient(rpcURL)
	if err != nil {
		logger.Fatal("创建以太坊客户端失败", errors.WithStack(err))
	}
	defer ethClient.Close()
	logger.Info("成功连接到Sepolia测试网")

	// 1. 区块查询流程
	blockService := block.NewService(ethClient)
	ctx := context.Background()

	logger.Info("开始查询区块信息", zap.Int64("target_block", cfg.Block.TargetNumber))
	targetBlock, err := blockService.GetBlockByNumber(ctx, cfg.Block.TargetNumber)
	if err != nil {
		logger.Fatal("区块查询失败", errors.WithStack(err))
	}
	blockService.PrintBlockInfo(targetBlock)

	// 2. 交易发送流程
	chainID := big.NewInt(cfg.Eth.ChainID)
	txService, err := transaction.NewService(
		ethClient,
		cfg.Transaction.PrivateKey,
		chainID,
	)
	if err != nil {
		logger.Fatal("创建交易服务失败", errors.WithStack(err))
	}

	logger.Info("开始发送转账交易",
		zap.String("recipient", cfg.Transaction.RecipientAddr.Hex()),
		zap.String("amount_wei", cfg.Transaction.AmountWei.String()))

	signedTx, err := txService.SendTransfer(
		ctx,
		cfg.Transaction.RecipientAddr,
		cfg.Transaction.AmountWei,
	)
	if err != nil {
		logger.Fatal("发送交易失败", errors.WithStack(err))
	}
	txService.PrintTxInfo(signedTx)

	// 等待退出信号，优雅关闭
	waitForShutdown()
}

// waitForShutdown 处理程序优雅退出
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("接收到退出信号，程序正在优雅关闭...")
}
