package config

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Config struct {
	// 链和网络配置
	InfuraAPIKey  string         `json:"infura_api_key"`
	PrivateKey    string         `json:"private_key"`
	RecipientAddr common.Address `json:"recipient_addr"`
	TransferWei   *big.Int       `json:"transfer_wei"`
	TargetBlock   *big.Int       `json:"target_block"`

	// 日志配置
	Env      string `json:"env"`       // 环境：development/production
	LogPath  string `json:"log_path"`  // 日志文件路径
	LogLevel string `json:"log_level"` // 日志级别：debug/info/warn/error
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Env:      "development",
		LogPath:  "./logs/sepolia-ops.log",
		LogLevel: "info",
		// 其他字段保持空，由用户在main中设置
	}
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if c.InfuraAPIKey == "" {
		return fmt.Errorf("Infura API Key 不能为空")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("私钥不能为空")
	}
	if c.RecipientAddr == (common.Address{}) {
		return fmt.Errorf("接收地址不能为空")
	}
	if c.TransferWei == nil || c.TransferWei.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("转账金额必须大于0")
	}
	if c.TargetBlock == nil || c.TargetBlock.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("区块号必须大于等于0")
	}
	return nil
}
