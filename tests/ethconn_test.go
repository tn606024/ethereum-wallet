package tests

import (
	"github.com/tn606024/ethwallet/conn"
	"github.com/tn606024/ethwallet/server"
	"net/http/httptest"
	"testing"
)

var ethconn *conn.EthConn



func TestServer(t *testing.T){
	ts := httptest.NewServer(server.SetupServer(TestNetwork, 8080))
	ethconn = conn.NewEthConn(ts.URL)
	defer ts.Close()

	_, err := ethconn.GetBalance(TestAddress)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetBlockNumber()
	if err != nil {t.Fatalf("%v\n", err)
	}

	_, err = ethconn.GetGasPrice()
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetNonce(TestAddress)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetTransaction(TestTransaction)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetNormalTransactions(TestAddress)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetInternalTransactions(TestAddress)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetTokenTransactions(TestAddress)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetEstimateGas(TestTransactionRequest)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = ethconn.GetErc20Balance(TestAddress)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}