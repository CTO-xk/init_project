package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config 应用全局配置
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	Chains   []ChainConfig  `yaml:"chains"`
	Points   PointsConfig   `yaml:"points"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// RabbitMQConfig MQ配置
type RabbitMQConfig struct {
	URL        string `yaml:"url"`
	Exchange   string `yaml:"exchange"`
	Queue      string `yaml:"queue"`
	RoutingKey string `yaml:"routing_key"`
}

// ChainConfig 区块链网络配置
type ChainConfig struct {
	Name            string `yaml:"name"`
	RPCURL          string `yaml:"rpc_url"`
	ChainID         int    `yaml:"chain_id"`
	ContractAddress string `yaml:"contract_address"`
	StartBlock      int64  `yaml:"start_block"`
	BlockDelay      int    `yaml:"block_delay"` // 区块确认延迟，固定为6
}

// PointsConfig 积分计算配置
type PointsConfig struct {
	Rate     float64 `yaml:"rate"`     // 积分比率，默认0.05
	Interval int     `yaml:"interval"` // 计算间隔（分钟），默认60
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	// 加载环境变量
	_ = godotenv.Load()

	// 读取YAML配置
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 从环境变量覆盖敏感配置
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		cfg.Database.Password = env
	}
	if env := os.Getenv("RABBITMQ_URL"); env != "" {
		cfg.RabbitMQ.URL = env
	}
	if env := os.Getenv("POINTS_RATE"); env != "" {
		if rate, err := strconv.ParseFloat(env, 64); err == nil {
			cfg.Points.Rate = rate
		}
	}

	// 强制设置区块延迟为6
	for i := range cfg.Chains {
		cfg.Chains[i].BlockDelay = 6
	}

	// 设置默认值
	if cfg.Points.Rate == 0 {
		cfg.Points.Rate = 0.05
	}
	if cfg.Points.Interval == 0 {
		cfg.Points.Interval = 60
	}

	return &cfg, nil
}
