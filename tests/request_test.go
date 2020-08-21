package tests

import (
	"fmt"
	"github.com/tn606024/ethwallet/ethclient"
	"github.com/tn606024/ethwallet/types"
	"testing"
)

var apiClient *ethclient.EthereumClient

func init() {
	path := types.LoadConfigPath()
	config, err := types.ImportConfig(path)
	if err != nil {
		fmt.Printf("Import Config occured error: %s", err)
	}
	apiClient = ethclient.NewEthereumClient(config.Ropsten.NodeUrl,config.Ropsten.EtherscanApiUrl, config.EtherscanApiKey, config.Network)
}

func TestEthereumClient_GetBlockNumber(t *testing.T){

	blocknum, err := apiClient.GetBlockNumber()
	if err != nil {
		t.Errorf("getBlockNumber failed: %s", err)
	}
	if blocknum == "" {
		t.Errorf("blockNumber is empty string")
	}
}

func TestEthereumClient_GetBalance(t *testing.T) {
	balance, err:= apiClient.GetBalance(TestAddress, types.Latest)
	if err != nil {
		t.Errorf("GetBalance failed: %s", err)
	}
	if balance == "" {
		t.Errorf("balance is empty string")
	}
}

func TestEthereumClient_GetErc20Balance(t *testing.T) {
	path := types.LoadConfigPath()
	config, _ := types.ImportConfig(path)
	_, err := apiClient.GetErc20ListBalance(config.Erc20List ,TestAddress)
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestEthereumClient_GetTransaction(t *testing.T) {
	_, err:= apiClient.GetTransaction(TestTransaction)
	if err != nil {
		t.Errorf("GetTransaction failed: %s", err)
	}
}
func TestEthereumClient_GetTransactionCount(t *testing.T) {
	nonce, err:= apiClient.GetTransactionCount(TestAddress, types.Latest)
	if err != nil {
		t.Errorf("GetTransactionCount failed: %s", err)
	}
	if nonce == "" {
		t.Errorf("nonce is empty string")
	}
}

func TestEthereumClient_GetGasPrice(t *testing.T) {
	gasPrice, err := apiClient.GetGasPrice()
	if err != nil{
		t.Errorf("GetGasPrice failed: %s", err)
	}
	if gasPrice == "" {
		t.Errorf("gasPrice is empty string")
	}
}

func TestEthereumClient_GetEstimateGas(t *testing.T) {
	for _ , test := range TestTransactionEstimateGas {
		esGas, err := apiClient.GetEstimateGas(&test.transaction)
		if err != nil {
			t.Errorf("GetEstimateGas failed: %s", err)

		}
		if esGas != test.gas {
			t.Errorf("the ans is %s, but we got %s", test.gas, esGas)
		}
	}
}

func TestEthereumClient_GetCall(t *testing.T) {
	for _ , test := range TestTranasctionCall {
		call, err := apiClient.GetCall(&test.transaction, types.Latest)
		if err != nil {
			t.Errorf("GetEstimateGas failed: %s", err)

		}
		if call != test.res {
			t.Errorf("the ans is %s, but we got %s", test.res, call)
		}
	}
}

func TestEthereumClient_GetNormalTransactions(t *testing.T) {
	_, err := apiClient.GetNormalTransactions(0,9999999, false, TestContractAddress)
	if err != nil {
		t.Errorf("GetNormalTransactions error: %s", err)
	}
}

func TestEthereumClient_GetErc20TokenTransactions(t *testing.T) {
	_, err := apiClient.GetErc20TokenTransactions(0,9999999, false, TestAddress)
	if err != nil {
		t.Errorf("GetErc20TokenTransactions error: %s", err)
	}
}

func TestEthereumClient_GetInternalTransactions(t *testing.T) {
	_ , err := apiClient.GetInternalTransactions(0,9999999, false, TestAddress)
	if err != nil {
		t.Errorf("GetInternalTransactions error: %s", err)
	}
}

func TestEthereumClient_GetLogs(t *testing.T) {
	topics := map[string]string{
		"topic0": "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"topic1": "0x00000000000000000000000051bf0b41Ba5B034f158CF1233f16bA5450F9355B",
		"topic1_2_opr": "or",
		"topic2": "0x00000000000000000000000051bf0b41Ba5B034f158CF1233f16bA5450F9355B",
	}
	_, err := apiClient.GetLogs(0,9999999, TestContractAddress, topics)
	if err != nil {
		t.Errorf("GetLogs error: %s", err)
	}
}