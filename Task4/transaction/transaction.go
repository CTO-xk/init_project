package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"sepolia-ops/config"
	"sepolia-ops/logger"
)

// TransactionService 交易处理服务
type TransactionService struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey // 发送方私钥
	chainID    *big.Int          // 链ID
}

// NewTransactionService 初始化交易服务
func NewTransactionService(client *ethclient.Client, privateKeyStr string, chainID *big.Int) (*TransactionService, error) {
	// 解析私钥（不含0x前缀）
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("私钥解析失败: %w", err)
	}

	logger.L().Info("交易服务初始化成功")
	return &TransactionService{
		client:     client,
		privateKey: privateKey,
		chainID:    chainID,
	}, nil
}

// CreateAndSendTransaction 创建并发送转账交易
func (s *TransactionService) CreateAndSendTransaction(
	ctx context.Context,
	recipient common.Address,
	amountWei *big.Int,
) (*types.Transaction, error) {
	// 1. 获取发送方地址
	publicKey := s.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥格式错误")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
	logger.L().Info("获取发送方地址", zap.String("address", fromAddr.Hex()))

	// 2. 获取Nonce（发送方下一个可用交易序号）
	nonce, err := s.client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %w", err)
	}
	logger.L().Info("获取交易Nonce", zap.Uint64("nonce", nonce))

	// 3. 获取建议的Gas价格
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取Gas价格失败: %w", err)
	}
	logger.L().Info("获取Gas价格", zap.String("gas_price_wei", gasPrice.String()))

	// 4. 构建未签名交易
	tx := types.NewTransaction(
		nonce,
		recipient,
		amountWei,
		config.DefaultGasLimit, // 默认Gas Limit（普通转账固定21000）
		gasPrice,
		nil, // 无附加数据
	)

	// 5. 签名交易
	signer := types.NewLondonSigner(s.chainID) // 兼容EIP-1559的签名器
	signedTx, err := types.SignTx(tx, signer, s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("交易签名失败: %w", err)
	}

	// 6. 发送交易到网络
	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("发送交易失败: %w", err)
	}

	logger.L().Info("交易发送成功",
		zap.String("tx_hash", signedTx.Hash().Hex()),
		zap.String("from", fromAddr.Hex()),
		zap.String("to", recipient.Hex()),
		zap.String("amount_wei", amountWei.String()))

	return signedTx, nil
}

// PrintTransactionInfo 打印交易信息
func (s *TransactionService) PrintTransactionInfo(tx *types.Transaction) {
	if tx == nil {
		fmt.Println("交易信息为空")
		return
	}
	fmt.Printf("\n===== 交易已发送 =====\n")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("查看详情: %s%s\n", config.SepoliaEtherscanURL, tx.Hash().Hex())
}
