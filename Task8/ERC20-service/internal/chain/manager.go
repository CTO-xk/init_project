package chain

import (
	"bytes"
	"context"
	"erc20-service/config"
	mq "erc20-service/internal/mq"
	"erc20-service/pkg/logger"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// Manager 多链事件监听管理器
type Manager struct {
	chains    []config.ChainConfig
	listeners map[string]*Listener
	abi       abi.ABI
	producer  *mq.PointsProducer
	exitChan  chan error
	wg        sync.WaitGroup
	log       *slog.Logger
}

// NewManager 创建管理器
func NewManager(chains []config.ChainConfig, abiBytes []byte, producer *mq.PointsProducer) *Manager {
	// 解析ABI
	contractABI, err := abi.JSON(bytes.NewReader(abiBytes))
	if err != nil {
		panic(fmt.Sprintf("解析ABI失败: %v", err))
	}

	return &Manager{
		chains:    chains,
		listeners: make(map[string]*Listener),
		abi:       contractABI,
		producer:  producer,
		exitChan:  make(chan error, 1),
		log:       logger.New("chain-manager"),
	}
}

// Start 启动所有链的监听器
func (m *Manager) Start(ctx context.Context) error {
	// 为每个链创建并启动监听器
	for _, chainCfg := range m.chains {
		listener, err := NewListener(chainCfg, m.abi, m.producer)
		if err != nil {
			return fmt.Errorf("创建监听器失败: %v", err)
		}
		m.listeners[chainCfg.Name] = listener

		m.wg.Add(1)
		go func(cfg config.ChainConfig, l *Listener) {
			defer m.wg.Done()
			m.log.Info("启动链监听器", "chain", cfg.Name)
			if err := l.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
				m.log.Error("链监听器异常退出", "chain", cfg.Name, "error", err)
				select {
				case m.exitChan <- fmt.Errorf("chain %s: %w", cfg.Name, err):
				default:
				}
			}
		}(chainCfg, listener)
	}

	m.wg.Wait()
	return nil
}

// GetChainNames 获取所有链名称
func (m *Manager) GetChainNames() []string {
	names := make([]string, 0, len(m.chains))
	for _, chain := range m.chains {
		names = append(names, chain.Name)
	}
	return names
}

// ExitChan 获取退出通知通道
func (m *Manager) ExitChan() <-chan error {
	return m.exitChan
}
