package config

import (
	"os"
	"strconv"
)

// Config 项目配置结构体
type Config struct {
	// 以太坊节点配置
	ETH RPCConfig `json:"eth"`
	// 合约配置
	Contract ContractConfig `json:"contract"`
	// 账户配置
	Account AccountConfig `json:"account"`
}

// RPCConfig 节点RPC配置
type RPCConfig struct {
	URL      string // 节点RPC地址
	Timeout  int    // 超时时间(秒)
	ChainID  int64  // 链ID (Sepolia测试网为11155111)
	GasPrice int64  // 默认GasPrice (gwei)
	GasLimit uint64 // 默认GasLimit
}

// ContractConfig 合约配置
type ContractConfig struct {
	Address string // 已部署的合约地址
}

// AccountConfig 账户配置
type AccountConfig struct {
	PrivateKey string // 用于交互的私钥
}

// Load 从环境变量加载配置
func Load() Config {
	// 读取链ID，默认Sepolia测试网
	chainID, _ := strconv.ParseInt(getEnv("ETH_CHAIN_ID", "11155111"), 10, 64)
	// 读取GasPrice，默认30gwei
	gasPrice, _ := strconv.ParseInt(getEnv("ETH_GAS_PRICE", "30"), 10, 64)
	// 读取GasLimit，默认300000
	gasLimit, _ := strconv.ParseUint(getEnv("ETH_GAS_LIMIT", "300000"), 10, 64)
	// 读取超时时间，默认30秒
	timeout, _ := strconv.Atoi(getEnv("ETH_TIMEOUT", "30"))

	return Config{
		ETH: RPCConfig{
			URL:      getEnv("ETH_RPC_URL", ""),
			Timeout:  timeout,
			ChainID:  chainID,
			GasPrice: gasPrice,
			GasLimit: gasLimit,
		},
		Contract: ContractConfig{
			Address: getEnv("COUNTER_CONTRACT_ADDRESS", ""),
		},
		Account: AccountConfig{
			PrivateKey: getEnv("ETH_PRIVATE_KEY", ""),
		},
	}
}

// 读取环境变量，带默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
