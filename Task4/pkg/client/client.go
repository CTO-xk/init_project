package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 以太坊客户端接口定义
type Client interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	PendingNonceAt(ctx context.Context, addr string) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	Close() error
}

// EthClient 基于ethclient的实现
type EthClient struct {
	client *ethclient.Client
}

// NewEthClient 创建客户端实例
func NewEthClient(rpcURL string) (*EthClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接RPC节点失败: %w", err)
	}
	return &EthClient{client: client}, nil
}

// BlockByNumber 获取指定区块
func (c *EthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.client.BlockByNumber(ctx, number)
}

// PendingNonceAt 获取指定地址的下一个可用nonce
func (c *EthClient) PendingNonceAt(ctx context.Context, addr string) (uint64, error) {
	return c.client.PendingNonceAt(ctx, common.HexToAddress(addr))
}

// SuggestGasPrice 获取当前区块的建议Gas价格
func (c *EthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.client.SuggestGasPrice(ctx)
}

// NetworkID 获取网络ID
func (c *EthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client.NetworkID(ctx)
}

// SendTransaction 发送交易
func (c *EthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.client.SendTransaction(ctx, tx)
}

// Close 关闭客户端
func (c *EthClient) Close() error {
	return nil
}
