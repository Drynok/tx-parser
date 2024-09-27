package parser

import (
	"context"
	"fmt"
	"sync"

	"github.com/Drynok/tx-parser/internal/model"
	rpc "github.com/Drynok/tx-parser/internal/rpc"
	"github.com/Drynok/tx-parser/internal/storage"
	"github.com/Drynok/tx-parser/pkg/logger"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []model.Transaction
}

type EthereumParser struct {
	CurrentBlock *model.Block
	mu           sync.RWMutex

	Storage   storage.Storage
	RPCClient rpc.Client
	Logger    logger.Logger
}

// NewEthereumParser construstor.
func NewEthereumParser(cli rpc.Client, storage storage.Storage, logger logger.Logger) *EthereumParser {
	return &EthereumParser{
		RPCClient: cli,
		Storage:   storage,
		Logger:    logger,
	}
}

// GetCurrentBlock returns the current block number.
func (p *EthereumParser) GetCurrentBlock() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.CurrentBlock.Number
}

func (p *EthereumParser) Subscribe(address string) bool {
	return p.Storage.Subscribe(address)
}

func (p *EthereumParser) GetTransactions(address string) []model.Transaction {
	return p.Storage.Transactions(address)
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
	latestBlock, err := p.RPCClient.GetLatestBlockNumber()
	if err != nil {
		p.Logger.Error("error fetching latest block:", err)
		return fmt.Errorf("error getting latest block number: %w", err)
	}

	if latestBlock <= p.GetCurrentBlock() {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, latestBlock-p.GetCurrentBlock())

	for int(p.CurrentBlock.Number) < latestBlock {
		wg.Add(1)
		go func(blockNumber int) {
			defer wg.Done()
			block, err := p.RPCClient.GetBlockByNumber(blockNumber)
			if err != nil {
				p.Logger.Error("failed to get block %d: %w", blockNumber, err)
				errChan <- fmt.Errorf("failed to get block %d: %w", blockNumber, err)
				return
			}
			if err := p.processBlock(block); err != nil {
				p.Logger.Error("error processing block", "error", err)
				errChan <- fmt.Errorf("error processing block %d: %w", blockNumber, err)
			}
		}(p.CurrentBlock.Number + 1)
		p.CurrentBlock.Number++
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		p.Logger.Error("Error during block processing", "error", err)
	}

	return nil
}

func (p *EthereumParser) processBlock(block *model.Block) error {
	for _, tx := range block.Transactions {
		if err := p.Storage.AddTransaction(tx.From, tx); err != nil {
			p.Logger.Error("error adding transaction: %w", err)
			return fmt.Errorf("error adding transaction: %w", err)
		}
		if err := p.Storage.AddTransaction(tx.To, tx); err != nil {
			p.Logger.Error("error adding transaction: %w", err)
			return fmt.Errorf("error adding transaction: %w", err)
		}
	}

	return nil
}

func (p *EthereumParser) isSubscribed(address string) bool {
	return p.Storage.IsSubscribed(address)
}
