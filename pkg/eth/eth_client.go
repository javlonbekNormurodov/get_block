package eth

import (
	"bytes"
	"encoding/json"
	"eth_tracker/internal/domain"
	"math/big"
	"net/http"
	"strconv"
)

type Client struct {
	rpcURL string
}

func NewEthClient(apiKey string) *Client {
	return &Client{rpcURL: "https://eth.getblock.io/mainnet/" + apiKey}
}

func (c *Client) call(method string, params []interface{}, result interface{}) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(c.rpcURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response struct {
		JSONRPC string          `json:"jsonrpc"`
		Result  json.RawMessage `json:"result"`
		Error   interface{}     `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	return json.Unmarshal(response.Result, result)
}

func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	var result string
	err := c.call("eth_blockNumber", []interface{}{}, &result)
	if err != nil {
		return nil, err
	}

	blockNumber := new(big.Int)
	blockNumber.SetString(result[2:], 16)
	return blockNumber, nil
}

func (c *Client) GetBlockByNumber(blockNumber int64) (*domain.Block, error) {
	var block domain.Block
	blockNumHex := "0x" + strconv.FormatInt(blockNumber, 16)
	err := c.call("eth_getBlockByNumber", []interface{}{blockNumHex, true}, &block)
	return &block, err
}

func (c *Client) GetTransactionByHash(hash string) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := c.call("eth_getTransactionByHash", []interface{}{hash}, &tx)
	return &tx, err
}
