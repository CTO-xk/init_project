package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/CTO-xk/init_project/Task4_02/config"
	bindings "github.com/CTO-xk/init_project/Task4_02/contracts/bindings"
	"github.com/CTO-xk/init_project/Task4_02/internal/counter"
	"github.com/CTO-xk/init_project/Task4_02/pkg/ethclient"
)

// 命令行参数：操作类型（部署/查询/增加/减少/重置）
var (
	action = flag.String("action", "query", "操作类型: deploy/query/increment/decrement/reset")
)

func main() {
	flag.Parse()

	// 1. 加载配置
	cfg := config.Load()
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 2. 初始化以太坊客户端
	client, err := ethclient.New(cfg.ETH)
	if err != nil {
		log.Fatalf("初始化以太坊客户端失败: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ 成功连接到Sepolia测试网")

	// 3. 处理不同操作
	switch *action {
	case "deploy":
		handleDeploy(client, cfg)
	case "query":
		handleQuery(client, cfg)
	case "increment":
		handleIncrement(client, cfg)
	case "decrement":
		handleDecrement(client, cfg)
	case "reset":
		handleReset(client, cfg)
	default:
		log.Fatalf("未知操作类型: %s，支持的操作: deploy/query/increment/decrement/reset", *action)
	}
}

// 验证配置有效性
func validateConfig(cfg config.Config) error {
	if cfg.ETH.URL == "" {
		return fmt.Errorf("请设置ETH_RPC_URL环境变量")
	}
	if *action != "deploy" && cfg.Contract.Address == "" {
		return fmt.Errorf("请设置COUNTER_CONTRACT_ADDRESS环境变量")
	}
	if cfg.Account.PrivateKey == "" {
		return fmt.Errorf("请设置ETH_PRIVATE_KEY环境变量")
	}
	return nil
}

// 处理合约部署
func handleDeploy(client *ethclient.Client, cfg config.Config) {
	// 创建临时服务（无需合约地址）
	// 注意：这里需要单独创建transactor，因为Service初始化依赖合约地址
	transactor, err := client.NewTransactor(cfg.Account.PrivateKey)
	if err != nil {
		log.Fatalf("创建交易签名器失败: %v", err)
	}

	// 部署合约（初始计数设为0）
	addr, tx, _, err := bindings.DeployCounter(transactor, client, big.NewInt(0))
	if err != nil {
		log.Fatalf("部署合约失败: %v", err)
	}

	fmt.Printf("📤 合约部署交易已发送: %s\n", tx.Hash().Hex())
	fmt.Printf("🔍 请等待确认，合约地址: %s\n", addr.Hex())
	fmt.Println("确认后可设置环境变量: export COUNTER_CONTRACT_ADDRESS=" + addr.Hex())
}

// 处理查询计数
func handleQuery(client *ethclient.Client, cfg config.Config) {
	service, err := counter.NewService(client, cfg)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	count, err := service.GetCount()
	if err != nil {
		log.Fatalf("查询计数失败: %v", err)
	}
	fmt.Printf("📊 合约 %s 当前计数: %d\n", service.Address(), count)
}

// 处理增加计数
func handleIncrement(client *ethclient.Client, cfg config.Config) {
	service, err := counter.NewService(client, cfg)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	txHash, err := service.Increment()
	if err != nil {
		log.Fatalf("增加计数失败: %v", err)
	}
	fmt.Printf("📈 增加计数成功，交易哈希: %s\n", txHash.Hex())

	// 再次查询确认
	count, _ := service.GetCount()
	fmt.Printf("📊 最新计数: %d\n", count)
}

// 处理减少计数
func handleDecrement(client *ethclient.Client, cfg config.Config) {
	service, err := counter.NewService(client, cfg)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	txHash, err := service.Decrement()
	if err != nil {
		log.Fatalf("减少计数失败: %v", err)
	}
	fmt.Printf("📉 减少计数成功，交易哈希: %s\n", txHash.Hex())

	// 再次查询确认
	count, _ := service.GetCount()
	fmt.Printf("📊 最新计数: %d\n", count)
}

// 处理重置计数
func handleReset(client *ethclient.Client, cfg config.Config) {
	service, err := counter.NewService(client, cfg)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	txHash, err := service.Reset()
	if err != nil {
		log.Fatalf("重置计数失败: %v", err)
	}
	fmt.Printf("🔄 重置计数成功，交易哈希: %s\n", txHash.Hex())

	// 再次查询确认
	count, _ := service.GetCount()
	fmt.Printf("📊 最新计数: %d\n", count)
}
