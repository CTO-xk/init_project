package counter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/CTO-xk/init_project/Task4_02/config" // 导入生成的合约绑定
	counter "github.com/CTO-xk/init_project/Task4_02/contracts/bindings"
	"github.com/CTO-xk/init_project/Task4_02/pkg/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Service 计数器合约服务
type Service struct {
	client     *ethclient.Client  // 以太坊客户端
	contract   *counter.Counter   // 合约实例
	address    common.Address     // 合约地址
	transactor *bind.TransactOpts // 交易签名器
}

// NewService 初始化服务
func NewService(client *ethclient.Client, cfg config.Config) (*Service, error) {
	// 解析合约地址
	contractAddr := common.HexToAddress(cfg.Contract.Address)

	// 创建合约实例
	contract, err := counter.NewCounter(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("创建合约实例失败: %w", err)
	}

	// 创建交易签名器
	transactor, err := client.NewTransactor(cfg.Account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("创建交易签名器失败: %w", err)
	}

	return &Service{
		client:     client,
		contract:   contract,
		address:    contractAddr,
		transactor: transactor,
	}, nil
}

// Deploy 部署新合约
func (s *Service) Deploy(initialCount int64) (common.Address, common.Hash, error) {
	// 部署合约
	addr, tx, _, err := counter.DeployCounter(
		s.transactor,
		s.client,
		big.NewInt(initialCount),
	)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("部署合约失败: %w", err)
	}
	return addr, tx.Hash(), nil
}

// GetCount 获取当前计数
func (s *Service) GetCount() (uint64, error) {
	count, err := s.contract.GetCount(&bind.CallOpts{
		Context: context.Background(),
	})
	if err != nil {
		return 0, fmt.Errorf("获取计数失败: %w", err)
	}
	return count.Uint64(), nil
}

// Increment 增加计数
func (s *Service) Increment() (common.Hash, error) {
	tx, err := s.contract.Increment(s.transactor)
	if err != nil {
		return common.Hash{}, fmt.Errorf("调用Increment失败: %w", err)
	}
	// 等待交易确认
	if err := s.client.WaitForTxConfirmation(tx.Hash()); err != nil {
		return tx.Hash(), fmt.Errorf("交易确认失败: %w", err)
	}
	return tx.Hash(), nil
}

// Decrement 减少计数
func (s *Service) Decrement() (common.Hash, error) {
	tx, err := s.contract.Decrement(s.transactor)
	if err != nil {
		return common.Hash{}, fmt.Errorf("调用Decrement失败: %w", err)
	}
	// 等待交易确认
	if err := s.client.WaitForTxConfirmation(tx.Hash()); err != nil {
		return tx.Hash(), fmt.Errorf("交易确认失败: %w", err)
	}
	return tx.Hash(), nil
}

// Reset 重置计数为0
func (s *Service) Reset() (common.Hash, error) {
	tx, err := s.contract.Reset(s.transactor)
	if err != nil {
		return common.Hash{}, fmt.Errorf("调用Reset失败: %w", err)
	}
	// 等待交易确认
	if err := s.client.WaitForTxConfirmation(tx.Hash()); err != nil {
		return tx.Hash(), fmt.Errorf("交易确认失败: %w", err)
	}
	return tx.Hash(), nil
}

// Address 返回合约地址
func (s *Service) Address() string {
	return s.address.Hex()
}
