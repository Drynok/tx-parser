package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Drynok/tx-parser/internal/model"
)

type Client interface {
	GetLatestBlockNumber() (int, error)
	GetBlockByNumber(blockNumber int) (*model.Block, error)
}

type RPCClient struct {
	url    string
	client *http.Client
}

func NewClient(url string) *RPCClient {
	return &RPCClient{
		url:    url,
		client: &http.Client{},
	}
}

type Request struct {
	ID      int         `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Response struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *RPCClient) GetLatestBlockNumber() (int, error) {
	resp, err := c.sendRequest("eth_blockNumber", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block number: %w", err)
	}

	var result string
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	blockNumber, err := strconv.ParseInt(result[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %w", err)
	}

	return int(blockNumber), nil
}

func (c *RPCClient) sendRequest(method string, params interface{}) (*Response, error) {
	request := Request{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(c.url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var jsonResp Response
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if jsonResp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %s (code: %d)", jsonResp.Error.Message, jsonResp.Error.Code)
	}

	return &jsonResp, nil
}

func (c *RPCClient) GetBlockByNumber(blockNumber int) (*model.Block, error) {
	params := []interface{}{fmt.Sprintf("0x%x", blockNumber), true}
	resp, err := c.sendRequest("eth_getBlockByNumber", params)
	if err != nil {
		return nil, err
	}

	var block model.Block
	if err := json.Unmarshal(resp.Result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	return &block, nil
}
