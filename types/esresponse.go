package types

import (
	"encoding/json"
)


type EsResponse struct {
	Status int `json:"status,string"`
	Message string `json:"message"`
	Result json.RawMessage `json:"result"`
}

type EsNormalTransaction struct {
	BlockNumber       int    `json:"blockNumber,string"`
	TimeStamp         Time   `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             int    `json:"nonce,string"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  int    `json:"transactionIndex,string"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             BigInt `json:"value"`
	Gas               int    `json:"gas,string"`
	GasPrice          BigInt `json:"gasPrice"`
	IsError           int    `json:"isError,string"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed int    `json:"cumulativeGasUsed,string"`
	GasUsed           int    `json:"gasUsed,string"`
	Confirmations     int    `json:"confirmations,string"`
}

func (t *EsNormalTransaction) String() (string, error){
	ts, err := json.MarshalIndent(t,"","	")
	if err != nil {
		return "", err
	}
	return string(ts) + "\n", nil
}

type EsInternalTansaction struct {
	Hash            string `json:"hash"`
	TraceID         string `json:"traceId"`
	BlockNumber     int    `json:"blockNumber,string"`
	TimeStamp       Time   `json:"timeStamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           BigInt `json:"value"`
	ContractAddress string `json:"contractAddress"`
	Input           string `json:"input"`
	Type            string `json:"type"`
	Gas             int    `json:"gas,string"`
	GasUsed         int    `json:"gasUsed,string"`
	IsError         int    `json:"isError,string"`
	ErrCode         string `json:"errCode"`
}

func (t *EsInternalTansaction) String() (string, error){
	ts, err := json.MarshalIndent(t,"","	")
	if err != nil {
		return "", err
	}
	return string(ts) + "\n", nil
}

type EsErc20TokenTransaction struct {
	BlockNumber       int    `json:"blockNumber,string"`
	TimeStamp         Time   `json:"timeStamp,string"`
	Hash              string `json:"hash"`
	Nonce             int    `json:"nonce,string"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             BigInt `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      int    `json:"tokenDecimal,string"`
	TransactionIndex  int    `json:"transactionIndex,string"`
	Gas               int    `json:"gas,string"`
	GasPrice          BigInt `json:"gasPrice"`
	GasUsed           int    `json:"gasUsed,string"`
	CumulativeGasUsed int    `json:"cumulativeGasUsed,string"`
	Input             string `json:"input"`
	Confirmations     int    `json:"confirmations,string"`
}

func (t *EsErc20TokenTransaction) String() (string, error){
	ts, err := json.MarshalIndent(t,"","	")
	if err != nil {
		return "", err
	}
	return string(ts) + "\n", nil
}

type EsLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TimeStamp        string   `json:"timeStamp"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}