package parser_test

import (
	"context"
	"testing"

	"github.com/Drynok/tx-parser/internal/model"
	"github.com/Drynok/tx-parser/internal/parser"
	rpc "github.com/Drynok/tx-parser/internal/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEthereumParser_NewEthereumParser(t *testing.T) {
	cli := &rpc.RPCClient{}
	storage := &mock.Storage{}
	logger := &mock.Logger{}

	parser := parser.NewEthereumParser(cli, storage, logger)

	assert.NotNil(t, parser)
	assert.Equal(t, cli, parser.rpcClient)
	assert.Equal(t, storage, parser.storage)
	assert.Equal(t, logger, parser.logger)
}

func TestEthereumParser_GetCurrentBlock(t *testing.T) {
	parser := &parser.EthereumParser{
		currentBlock: &model.Block{Number: 10},
	}

	assert.Equal(t, 10, parser.GetCurrentBlock())
}

func TestEthereumParser_Subscribe(t *testing.T) {
	storage := &mock.Storage{}
	storage.On("Subscribe", "address").Return(true)

	parser := &parser.EthereumParser{
		storage: storage,
	}

	assert.True(t, parser.Subscribe("address"))
}

func TestEthereumParser_GetTransactions(t *testing.T) {
	storage := &mock.Storage{}
	storage.On("Transactions", "address").Return([]model.Transaction{{From: "from", To: "to"}})

	parser := &parser.EthereumParser{
		storage: storage,
	}

	transactions := parser.GetTransactions("address")
	assert.Len(t, transactions, 1)
	assert.Equal(t, "from", transactions[0].From)
	assert.Equal(t, "to", transactions[0].To)
}

func TestEthereumParser_Start(t *testing.T) {
	cli := &mock.RPCClient{}
	storage := &mock.Storage{}
	logger := &mock.Logger{}

	parser := &parser.EthereumParser{
		rpcClient: cli,
		storage:   storage,
		logger:    logger,
	}

	ctx := context.Background()
	err := parser.Start(ctx)
	assert.NoError(t, err)
}

func TestEthereumParser_pollBlocks(t *testing.T) {
	cli := &mock.RPCClient{}
	storage := &mock.Storage{}
	logger := &mock.Logger{}

	parser := &parser.EthereumParser{
		rpcClient: cli,
		storage:   storage,
		logger:    logger,
	}

	err := parser.pollBlocks()
	assert.NoError(t, err)
}

func TestEthereumParser_processBlock(t *testing.T) {
	storage := &mock.Storage{}
	storage.On("AddTransaction", "from", mock.Anything).Return(nil)
	storage.On("AddTransaction", "to", mock.Anything).Return(nil)

	parser := &parser.EthereumParser{
		storage: storage,
	}

	block := &model.Block{
		Transactions: []model.Transaction{{From: "from", To: "to"}},
	}

	err := parser.processBlock(block)
	assert.NoError(t, err)
}

func TestEthereumParser_isSubscribed(t *testing.T) {
	storage := &mock.Storage{}
	storage.On("IsSubscribed", "address").Return(true)

	parser := &parser.EthereumParser{
		storage: storage,
	}

	assert.True(t, parser.isSubscribed("address"))
}
