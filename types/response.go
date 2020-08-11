package types

import (
	"encoding/json"
)

type Response struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
}

type ResponseError struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Error   EthereumError `json:"error"`
}

type EthereumError struct {
	Code 		int		`json:"code"`
	Message		string	`json:"message"`
}

type NodeTransaction struct {
	BlockNumber      IntHex    `json:"blockNumber"`
	Hash             string    `json:"hash"`
	Nonce            IntHex    `json:"nonce"`
	BlockHash        string    `json:"blockHash"`
	TransactionIndex IntHex    `json:"transactionIndex"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	Value            BigIntHex `json:"value"`
	Gas              IntHex    `json:"gas"`
	GasPrice         BigIntHex `json:"gasPrice"`
	Input            string    `json:"input"`
}