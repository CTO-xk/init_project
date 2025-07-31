package block

import (
	"github.com/CTO-xk/init_project/Task4/logger"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
)

// BlockService 区块查询服务
type BlockService struct {
	client *ethclient.Client
}

// NewBlockService 初始化区块服务
func NewBlockService(client *ethclient.Client) *BlockService {
	return &BlockService{client: client}
}

// GetBlockByNumber 查询指定区块号的区块信息
func (s *BlockService) GetBlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	logger.L().Info("开始查询区块", zap.Int64("block_number", blockNumber.Int64()))
	block, err := s.client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		logger.L().Error("获取区块失败",
			zap.Int64("block_number", blockNumber.Int64()),
			zap.Error(err)
		return nil, fmt.Errorf("获取区块失败: %w", err)
	}
	logger.L().Info("区块查询成功",
		zap.Int64("block_number", blockNumber.Int64()),
		zap.String("block_hash", block.Hash().Hex()))

	return block, nil
}

// PrintBlockInfo 打印区块详细信息到控制台
func (s *BlockService) PrintBlockInfo(block *types.Block) {
	if block == nil {
		fmt.Println("区块信息为空")
		return
	}
	fmt.Printf("\n===== 区块信息 (区块号 #%d) =====\n", block.Number().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())
	fmt.Printf("区块大小: %d bytes\n", block.Size())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
}
