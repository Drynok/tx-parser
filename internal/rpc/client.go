package rpc

import (
	"github.com/Drynok/tx-parser/internal/model"
)

type Client interface {
	GetLatestBlockNumber() (int, error)
	GetBlockByNumber(blockNumber int) (*model.Block, error)
}
