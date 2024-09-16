package parser

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Drynok/tx-parser/internal/model"
	rpc "github.com/Drynok/tx-parser/internal/rpc"
	"github.com/Drynok/tx-parser/internal/storage"
	"github.com/Drynok/tx-parser/pkg/logger"
)

type EthereumParser struct {
	currentBlock *model.Block
	subscribers  map[string]bool
	transactions map[string][]model.Transaction

	mu sync.RWMutex

	// TODO: Move all transactions and subscribers to storage
	storage   storage.Storage
	rpcClient rpc.Client
	logger    logger.Logger
}

func NewEthereumParser(cli rpc.Client, storage storage.Storage, logger logger.Logger) *EthereumParser {
	return &EthereumParser{
		storage:   storage,
		rpcClient: cli,
		logger:    logger,
	}
}

func (p *EthereumParser) GetCurrentBlock() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentBlock.Number
}

func (p *EthereumParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.subscribers[address]; exists {
		return false
	}
	p.subscribers[address] = true
	return true
}

func (p *EthereumParser) GetTransactions(address string) []model.Transaction {
	return p.storage.Transactions(address)
}

func (p *EthereumParser) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := p.pollBlocks(); err != nil {
				return fmt.Errorf("error processing blocks: %w", err)
			}
		}
	}
}

func (p *EthereumParser) pollBlocks() error {
	latestBlock, err := p.rpcClient.GetLatestBlockNumber()
	if err != nil {
		log.Println("Error fetching latest block:", err)
		return fmt.Errorf("error getting latest block number: %w", err)
	}

	if latestBlock <= p.GetCurrentBlock() {
		return nil
	}

	for int(p.currentBlock.Number) < latestBlock {
		p.currentBlock.Number++
		block, err := p.rpcClient.GetBlockByNumber(p.currentBlock.Number)
		if err != nil {
			p.logger.Error("Failed to get block", "block", p.currentBlock, "error", err)
			continue
		}

		if err := p.processBlock(block); err != nil {
			return fmt.Errorf("error processing block %d: %w", block, err)
		}

		p.mu.Lock()
		p.currentBlock = block
		p.mu.Unlock()
	}

	return nil
}

func (p *EthereumParser) processBlock(block *model.Block) error {
	for _, tx := range block.Transactions {
		if err := p.storage.AddTransaction(tx.From, tx); err != nil {
			p.logger.Error("error adding transaction: %w", err)
			return fmt.Errorf("error adding transaction: %w", err)
		}
		if err := p.storage.AddTransaction(tx.To, tx); err != nil {
			p.logger.Error("error adding transaction: %w", err)
			return fmt.Errorf("error adding transaction: %w", err)
		}
	}
	return nil
}

func (p *EthereumParser) isSubscribed(address string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.subscribers[address]
	return ok
}
