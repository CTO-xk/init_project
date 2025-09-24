package chain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"erc20-service/config"
	"erc20-service/internal/db"
	mq "erc20-service/internal/mq"
	"erc20-service/pkg/logger"
	"log/slog"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Listener 单链事件监听器
type Listener struct {
	chainCfg     config.ChainConfig
	client       *ethclient.Client
	contractABI  abi.ABI
	contractAddr common.Address
	producer     *mq.PointsProducer
	lastBlock    int64
	log          *slog.Logger
}

// NewListener 创建监听器
func NewListener(cfg config.ChainConfig, contractABI abi.ABI, producer *mq.PointsProducer) (*Listener, error) {
	// 连接RPC节点
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("连接RPC失败: %v", err)
	}

	// 获取上次处理的区块号
	lastBlock, err := db.GetLastProcessedBlock(cfg.Name)
	if err != nil || lastBlock == 0 {
		lastBlock = uint64(cfg.StartBlock)
	}

	return &Listener{
		chainCfg:     cfg,
		client:       client,
		contractABI:  contractABI,
		contractAddr: common.HexToAddress(cfg.ContractAddress),
		producer:     producer,
		lastBlock:    int64(lastBlock),
		log:          logger.New(fmt.Sprintf("chain:%s", cfg.Name)),
	}, nil
}

// Start 启动监听器
func (l *Listener) Start(ctx context.Context) error {
	l.log.Info("启动监听器",
		"contract", l.contractAddr.Hex(),
		"start_block", l.lastBlock,
		"block_delay", l.chainCfg.BlockDelay,
	)

	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		l.client.Close()
		l.log.Info("监听器已停止")
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := l.processBlocks(ctx); err != nil {
				l.log.Error("处理区块失败", "error", err)
			}
		}
	}
}

// 处理区块范围
func (l *Listener) processBlocks(ctx context.Context) error {
	// 获取最新区块
	latestBlock, err := l.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("获取最新区块失败: %v", err)
	}

	// 计算目标区块（应用6区块延迟）
	targetBlock := int64(latestBlock) - int64(l.chainCfg.BlockDelay)
	if targetBlock <= l.lastBlock {
		// 减少调试日志频率，避免闪烁
		return nil
	}

	l.log.Info("开始处理区块", "from", l.lastBlock+1, "to", targetBlock)

	// 过滤事件
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(l.lastBlock + 1),
		ToBlock:   big.NewInt(targetBlock),
		Addresses: []common.Address{l.contractAddr},
	}

	logs, err := l.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("过滤日志失败: %v", err)
	}

	// 处理所有事件
	for _, vLog := range logs {
		if err := l.processLog(vLog); err != nil {
			l.log.Warn("处理日志失败", "tx_hash", vLog.TxHash.Hex(), "error", err)
		}
	}

	// 更新最后处理的区块
	if err := db.UpdateLastProcessedBlock(l.chainCfg.Name, uint64(targetBlock)); err != nil {
		return fmt.Errorf("更新区块号失败: %v", err)
	}
	l.lastBlock = targetBlock
	l.log.Info("区块处理完成", "last_block", targetBlock)

	return nil
}

// 处理单条日志
func (l *Listener) processLog(vLog types.Log) error {
	event, err := l.contractABI.EventByID(vLog.Topics[0])
	if err != nil {
		return fmt.Errorf("未知事件ID: %v", err)
	}

	switch event.Name {
	case "Mint":
		return l.handleMint(vLog)
	case "Burn":
		return l.handleBurn(vLog)
	case "Transfer":
		return l.handleTransfer(vLog)
	default:
		l.log.Debug("忽略未知事件", "name", event.Name)
		return nil
	}
}

// 处理Mint事件
func (l *Listener) handleMint(vLog types.Log) error {
	var event struct {
		To        common.Address
		Amount    *big.Int
		Timestamp uint64
	}
	if err := l.contractABI.UnpackIntoInterface(&event, "Mint", vLog.Data); err != nil {
		return fmt.Errorf("解析Mint事件失败: %v", err)
	}

	// 计算新余额
	newBalance := l.calculateNewBalance(event.To.Hex(), event.Amount, true)

	// 记录余额变动
	change := db.BalanceChange{
		ChainName:    l.chainCfg.Name,
		UserAddress:  event.To.Hex(),
		EventType:    "mint",
		Amount:       event.Amount.String(),
		BalanceAfter: newBalance,
		BlockNumber:  vLog.BlockNumber,
		EventTime:    time.Unix(int64(event.Timestamp), 0),
		TxHash:       vLog.TxHash.Hex(),
	}

	if err := db.RecordBalanceChange(change); err != nil {
		return fmt.Errorf("记录Mint事件失败: %v", err)
	}

	l.log.Info("处理Mint事件",
		"to", event.To.Hex(),
		"amount", event.Amount.String(),
		"tx", vLog.TxHash.Hex(),
	)
	return nil
}

// 处理Burn事件
func (l *Listener) handleBurn(vLog types.Log) error {
	var event struct {
		From      common.Address
		Amount    *big.Int
		Timestamp uint64
	}
	if err := l.contractABI.UnpackIntoInterface(&event, "Burn", vLog.Data); err != nil {
		return fmt.Errorf("解析Burn事件失败: %v", err)
	}

	// 计算新余额
	newBalance := l.calculateNewBalance(event.From.Hex(), event.Amount, false)

	// 记录余额变动
	change := db.BalanceChange{
		ChainName:    l.chainCfg.Name,
		UserAddress:  event.From.Hex(),
		EventType:    "burn",
		Amount:       event.Amount.String(),
		BalanceAfter: newBalance,
		BlockNumber:  vLog.BlockNumber,
		EventTime:    time.Unix(int64(event.Timestamp), 0),
		TxHash:       vLog.TxHash.Hex(),
	}

	if err := db.RecordBalanceChange(change); err != nil {
		return fmt.Errorf("记录Burn事件失败: %v", err)
	}

	l.log.Info("处理Burn事件",
		"from", event.From.Hex(),
		"amount", event.Amount.String(),
		"tx", vLog.TxHash.Hex(),
	)
	return nil
}

// 处理Transfer事件
func (l *Listener) handleTransfer(vLog types.Log) error {
	if len(vLog.Topics) < 3 {
		return fmt.Errorf("Transfer事件参数不足")
	}

	from := common.HexToAddress(vLog.Topics[1].Hex())
	to := common.HexToAddress(vLog.Topics[2].Hex())
	amount := new(big.Int).SetBytes(vLog.Data)
	eventTime := time.Now()

	zero := common.Address{}

	// 铸造：from 为零地址
	if from == zero && to != zero {
		newBalance := l.calculateNewBalance(to.Hex(), amount, true)
		change := db.BalanceChange{
			ChainName:    l.chainCfg.Name,
			UserAddress:  to.Hex(),
			EventType:    "mint",
			Amount:       amount.String(),
			BalanceAfter: newBalance,
			BlockNumber:  vLog.BlockNumber,
			EventTime:    eventTime,
			TxHash:       vLog.TxHash.Hex(),
		}
		if err := db.RecordBalanceChange(change); err != nil {
			return fmt.Errorf("记录铸造事件失败: %v", err)
		}
		l.log.Info("处理Mint(零地址)事件", "to", to.Hex(), "amount", amount.String(), "tx", vLog.TxHash.Hex())
		return nil
	}

	// 销毁：to 为零地址
	if to == zero && from != zero {
		newBalance := l.calculateNewBalance(from.Hex(), amount, false)
		change := db.BalanceChange{
			ChainName:    l.chainCfg.Name,
			UserAddress:  from.Hex(),
			EventType:    "burn",
			Amount:       amount.String(),
			BalanceAfter: newBalance,
			BlockNumber:  vLog.BlockNumber,
			EventTime:    eventTime,
			TxHash:       vLog.TxHash.Hex(),
		}
		if err := db.RecordBalanceChange(change); err != nil {
			return fmt.Errorf("记录销毁事件失败: %v", err)
		}
		l.log.Info("处理Burn(零地址)事件", "from", from.Hex(), "amount", amount.String(), "tx", vLog.TxHash.Hex())
		return nil
	}

	// 处理转出
	if from != zero {
		newBalance := l.calculateNewBalance(from.Hex(), amount, false)
		change := db.BalanceChange{
			ChainName:    l.chainCfg.Name,
			UserAddress:  from.Hex(),
			EventType:    "transfer_out",
			Amount:       amount.String(),
			BalanceAfter: newBalance,
			BlockNumber:  vLog.BlockNumber,
			EventTime:    eventTime,
			TxHash:       vLog.TxHash.Hex(),
		}
		if err := db.RecordBalanceChange(change); err != nil {
			return fmt.Errorf("记录转出事件失败: %v", err)
		}
	}

	// 处理转入
	if to != zero {
		newBalance := l.calculateNewBalance(to.Hex(), amount, true)
		change := db.BalanceChange{
			ChainName:    l.chainCfg.Name,
			UserAddress:  to.Hex(),
			EventType:    "transfer_in",
			Amount:       amount.String(),
			BalanceAfter: newBalance,
			BlockNumber:  vLog.BlockNumber,
			EventTime:    eventTime,
			TxHash:       vLog.TxHash.Hex(),
		}
		if err := db.RecordBalanceChange(change); err != nil {
			return fmt.Errorf("记录转入事件失败: %v", err)
		}
	}

	l.log.Info("处理Transfer事件",
		"from", from.Hex(),
		"to", to.Hex(),
		"amount", amount.String(),
		"tx", vLog.TxHash.Hex(),
	)
	return nil
}

// 计算新余额
func (l *Listener) calculateNewBalance(userAddr string, amount *big.Int, isIncrease bool) string {
	current, err := db.GetUserCurrentBalance(l.chainCfg.Name, userAddr)
	if err != nil {
		current = "0"
	}

	currentBig := new(big.Int)
	currentBig.SetString(current, 10)

	if isIncrease {
		currentBig.Add(currentBig, amount)
	} else {
		currentBig.Sub(currentBig, amount)
	}

	return currentBig.String()
}
