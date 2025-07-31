package config

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// 常量定义
const (
	DefaultGasLimit     = 21000                                                        // 普通转账默认Gas Limit
	SepoliaEtherscanURL = "https://eth-mainnet.g.alchemy.com/v2/XE2nmOMCIb6XVkP4Rj7Ar" // 区块浏览器URL
)

// Config 应用总配置
type Config struct {
	Eth         EthConfig         `mapstructure:"eth"`
	Log         LogConfig         `mapstructure:"log"`
	Block       BlockConfig       `mapstructure:"block"`
	Transaction TransactionConfig `mapstructure:"transaction"`
}

// EthConfig 以太坊网络配置
type EthConfig struct {
	InfuraAPIKey string `mapstructure:"infura_api_key" env:"INFURA_API_KEY"`
	ChainID      int64  `mapstructure:"chain_id"`
}

// LogConfig 日志配置
type LogConfig struct {
	Env        string `mapstructure:"env"`
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// BlockConfig 区块查询配置
type BlockConfig struct {
	TargetNumber int64 `mapstructure:"target_number"`
}

// TransactionConfig 交易配置
type TransactionConfig struct {
	PrivateKey    string         `mapstructure:"private_key" env:"PRIVATE_KEY"`
	RecipientAddr common.Address `mapstructure:"recipient_addr"`
	AmountWei     *big.Int       `mapstructure:"amount_wei"`
}

// Load 从配置文件加载配置
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		// 默认配置文件路径
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取工作目录失败: %w", err)
		}
		// 检查是否在cmd/main目录下运行
		if strings.HasSuffix(wd, "cmd/main") {
			configPath = filepath.Join(wd, "..", "..", "pkg", "config", "config.yaml")
		} else {
			configPath = filepath.Join(wd, "pkg", "config", "config.yaml")
		}
	}

	// 配置viper
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()                                   // 允许环境变量覆盖
	viper.SetEnvPrefix("SEPOLIA")                          // 环境变量前缀
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 替换分隔符

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	var cfg Config
	decoderConfig := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &cfg,
		TagName:  "mapstructure", // 指定结构体标签
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			// 处理字符串到big.Int的转换
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				if f.Kind() == reflect.String && t == reflect.TypeOf((*big.Int)(nil)) {
					str := data.(string)
					if str == "" {
						return big.NewInt(0), nil
					}
					val, ok := new(big.Int).SetString(str, 10)
					if !ok {
						return nil, fmt.Errorf("无法解析big.Int: %s", str)
					}
					return val, nil
				}
				if f.Kind() == reflect.Int64 && t == reflect.TypeOf((*big.Int)(nil)) {
					val := data.(int64)
					return big.NewInt(val), nil
				}
				if f.Kind() == reflect.Int && t == reflect.TypeOf((*big.Int)(nil)) {
					val := data.(int)
					return big.NewInt(int64(val)), nil
				}
				return data, nil
			},
			// 处理字符串到common.Address的转换
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				if f.Kind() == reflect.String && t == reflect.TypeOf(common.Address{}) {
					str := data.(string)
					if str == "" {
						return common.Address{}, nil
					}
					if !common.IsHexAddress(str) {
						return nil, fmt.Errorf("无效的以太坊地址: %s", str)
					}
					return common.HexToAddress(str), nil
				}
				return data, nil
			},
		),
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, fmt.Errorf("创建配置解码器失败: %w", err)
	}

	if err := decoder.Decode(viper.AllSettings()); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	// 验证配置有效性
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Eth.InfuraAPIKey == "" {
		return fmt.Errorf("eth.infura_api_key不能为空")
	}
	if c.Eth.ChainID <= 0 {
		return fmt.Errorf("eth.chain_id必须大于0")
	}
	if c.Transaction.PrivateKey == "" {
		return fmt.Errorf("transaction.private_key不能为空，请在配置文件中设置私钥或通过环境变量SEPOLIA_PRIVATE_KEY设置")
	}
	if c.Transaction.RecipientAddr == (common.Address{}) {
		return fmt.Errorf("transaction.recipient_addr不能为空")
	}
	if c.Transaction.AmountWei == nil || c.Transaction.AmountWei.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("transaction.amount_wei必须大于0")
	}
	if c.Block.TargetNumber < 0 {
		return fmt.Errorf("block.target_number不能为负数")
	}
	return nil
}
