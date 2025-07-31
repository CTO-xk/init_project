package client

import (
	"context"
	"fmt"
	"github.com/CTO-xk/init_project/Task4/logger"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
)

// ETHClient 以太坊客户端封装
type ETHClient struct {
	client *ethclient.Client
}

// NewETHClient 初始化客户端（连接到Sepolia测试网）
func NewETHClient(infuraAPIKey string) (*ETHClient, error) {
	rpcURL := "https://sepolia.infura.io/v3/" + infuraAPIKey
	logger.L().Info("连接到以太坊节点", zap.String("url", rpcURL))
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接到以太坊节点失败: %w", err)
	}
	logger.L().Info("成功连接到以太坊节点")
	return &ETHClient{client: client}, nil
}

// GetClient 获取底层ethclient实例
func (c *ETHClient) GetClient() *ethclient.Client {
	return c.client
}

// Close 关闭客户端连接
func (c *ETHClient) Close() {
	// ethclient无显式Close方法，底层连接会自动管理
}

// GetChainID 获取当前网络链ID
func (c *ETHClient) GetChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取链ID失败: %w", err)
	}
	logger.L().Info("获取链ID成功", zap.Uint64("chain_id", chainID.Uint64()))
	return chainID, nil
}
