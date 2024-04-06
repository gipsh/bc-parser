//go:generate mockgen -source=eth.go -destination=mocks/eth.go -package=mock

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Path: internal/observer/eth_observer.go
// Implement JSON-RPC client for Ethereum node
// Two methods are implemented: GetBlockNumber and GetBlockByNumber

type Client interface {
	GetBlockNumber() (string, error)
	GetBlockByNumber(blockNumber string) (*Block, error)
}

type EthClient struct {
	url string
}

func NewEthClient(url string) *EthClient {
	return &EthClient{url: url}
}

// curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
func (c *EthClient) GetBlockNumber() (string, error) {

	url := "https://cloudflare-eth.com"
	method := "POST"

	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var response JsonRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if response.Error != nil {
		fmt.Println(response.Error.Message)
		return "", fmt.Errorf(response.Error.Message)
	}

	var blockNumber string
	err = json.Unmarshal(response.Result, &blockNumber)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return blockNumber, nil

}

// curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}'
func (c *EthClient) GetBlockByNumber(blockNumber string) (*Block, error) {

	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":1}`, blockNumber))

	client := &http.Client{}
	req, err := http.NewRequest(method, c.url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var response JsonRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if response.Error != nil {
		fmt.Println(response.Error.Message)
		return nil, fmt.Errorf(response.Error.Message)
	}

	var block Block
	err = json.Unmarshal(response.Result, &block)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &block, nil
}
