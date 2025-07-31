package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/CTO-xk/init_project/Task4/pkg/client"
	"github.com/CTO-xk/init_project/Task4/pkg/config"
	"github.com/CTO-xk/init_project/Task4/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"go.uber.org/zap"

)

// Service 交易服务
type Service struct {
	client     client.Client
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}

// NewService 创建交易服务实例
func NewService(c client.Client, privateKeyStr string, chainID *big.Int) (*Service, error) {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, errors.Wrap(err, "解析私钥失败")
	}

	return &Service{
		client:     c,
		privateKey: privateKey,
		chainID:    chainID,
	}, nil
}

// SendTransfer 发送以太币转账交易
func (s *Service) SendTransfer(
	ctx context.Context,
	to common.Address,
	amount *big.Int,
) (*types.Transaction, error) {
	// 获取发送方地址
	publicKey := s.privateKey.Public().(*ecdsa.PublicKey)
	fromAddr := crypto.PubkeyToAddress(*publicKey)

	// 获取nonce
	nonce, err := s.client.PendingNonceAt(ctx, fromAddr.Hex())
	if err != nil {
		return nil, errors.Wrap(err, "获取nonce失败")
	}

	// 获取建议的gas价格
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取gas价格失败")
	}

	// 构建未签名交易
	tx := types.NewTransaction(
		nonce,
		to,
		amount,
		config.DefaultGasLimit, // 使用默认gas限制
		gasPrice,
		nil, // 无附加数据
	)

	// 签名交易
	signer := types.NewLondonSigner(s.chainID)
	signedTx, err := types.SignTx(tx, signer, s.privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "交易签名失败")
	}

	// 发送交易到网络
	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return nil, errors.Wrap(err, "发送交易失败")
	}

	logger.Info("交易发送成功", 
		zap.String("tx_hash", signedTx.Hash().Hex()),
		zap.String("from", fromAddr.Hex()),
		zap.String("to", to.Hex()))
	
	return signedTx, nil
}

// PrintTxInfo 打印交易信息到控制台
func (s *Service) PrintTxInfo(tx *types.Transaction) {
	if tx == nil {
		return
	}

	// 转换为ETH显示
	ethAmount := new(big.Float).Quo(
		new(big.Float).SetInt(tx.Value()),
		big.NewFloat(1e18),
	)

	fmt.Println("\n===== 交易信息 =====")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("发送方Nonce: %d\n", tx.Nonce())
	fmt.Printf("接收地址: %s\n", tx.To().Hex())
	fmt.Printf("转账金额: %.6f ETH\n", ethAmount)
	fmt.Printf("Gas价格: %s Wei\n", tx.GasPrice().String())
	fmt.Printf("Gas限制: %d\n", tx.Gas())
	fmt.Printf("区块浏览器查看: %s%s\n", config.SepoliaEtherscanURL, tx.Hash().Hex())
	fmt.Println("====================")
}
