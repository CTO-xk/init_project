package cmd

import (
	"erc20-service/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "erc20-tracker",
		Short: "多链ERC20事件追踪与积分计算服务",
		Long:  "支持Sepolia/Base Sepolia链事件监听、余额重建、MQ积分计算",
	}
	log = logger.New("root-cmd")
)

// Execute 执行根命令
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error("服务启动失败", "error", err)
		os.Exit(1)
	}
}
func init() {
	// 全局flag：配置文件路径
	RootCmd.PersistentFlags().StringP("config", "c", "config/config.yaml", "配置文件路径")
}
