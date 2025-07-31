package block

import (
	"context"
	"fmt"
	"math/big"

	"github.com/CTO-xk/init_project/Task4/pkg/client"
	"github.com/CTO-xk/init_project/Task4/pkg/logger"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// Service 区块服务
type Service struct {
	client client.Client // 依赖接口而非具体实现
}

// NewService 创建区块服务实例
func NewService(c client.Client) *Service {
	return &Service{client: c}
}

// GetBlockByNumber 查询指定区块号的区块信息
func (s *Service) GetBlockByNumber(ctx context.Context, blockNumber int64) (*types.Block, error) {
	num := big.NewInt(blockNumber)
	block, err := s.client.BlockByNumber(ctx, num)
	if err != nil {
		return nil, errors.Wrapf(err, "查询区块#%d失败", blockNumber)
	}
	return block, nil
}

// PrintBlockInfo 打印区块详细信息到控制台
func (s *Service) PrintBlockInfo(block *types.Block) {
	if block == nil {
		logger.Warn("区块信息为空")
		return
	}

	fmt.Println("\n===== 区块信息 =====")
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
	fmt.Printf("区块大小: %d bytes\n", block.Size())
	fmt.Println("====================")
}
