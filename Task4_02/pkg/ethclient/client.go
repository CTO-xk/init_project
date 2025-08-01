package ethclient

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/CTO-xk/init_project/Task4_02/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 以太坊客户端封装
type Client struct {
	*ethclient.Client
	chainID *big.Int
	config  config.RPCConfig
}

// New 初始化以太坊客户端
func New(cfg config.RPCConfig) (*Client, error) {
	// 连接节点
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("连接节点失败: %w", err)
	}

	// 验证链ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取链ID失败: %w", err)
	}
	if chainID.Int64() != cfg.ChainID {
		return nil, fmt.Errorf("链ID不匹配: 配置=%d, 实际=%d", cfg.ChainID, chainID.Int64())
	}

	return &Client{
		Client:  client,
		chainID: chainID,
		config:  cfg,
	}, nil
}

// NewTransactor 创建带签名的交易器
func (c *Client) NewTransactor(privateKey string) (*bind.TransactOpts, error) {
	// 解析私钥
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	// 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(key, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("创建交易器失败: %w", err)
	}

	// 设置Gas参数
	auth.GasLimit = c.config.GasLimit
	auth.GasPrice = big.NewInt(c.config.GasPrice * 1e9) // 转换为wei (1gwei = 1e9 wei)
	auth.Context = c.newContext()

	return auth, nil
}

// WaitForTxConfirmation 等待交易确认
func (c *Client) WaitForTxConfirmation(txHash common.Hash) error {
	ctx := c.newContext()
	// 轮询交易状态
	for {
		receipt, err := c.TransactionReceipt(ctx, txHash)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				// 交易未确认，继续等待
				time.Sleep(3 * time.Second)
				continue
			}
			return fmt.Errorf("获取交易收据失败: %w", err)
		}

		// 交易已确认（状态码0表示成功）
		if receipt.Status == 1 {
			return nil
		}
		return fmt.Errorf("交易执行失败，状态码: %d", receipt.Status)
	}
}

// 创建带超时的上下文
func (c *Client) newContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(c.config.Timeout)*time.Second)
	return ctx
}
