package conn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ian/ethwallet/types"
	"github.com/ian/ethwallet/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

type EthConn struct {
	conn	*http.Client
	url		string
}

type Response struct {
	Result json.RawMessage `json:"result"`
}

func NewEthConn(url string) *EthConn{
	if url == "" {
		url = "http://127.0.0.1:8080"
	}
	return &EthConn{
		conn: &http.Client{
			Timeout: 30 * time.Second,
		},
		url: url,
	}
}

func (c *EthConn) get(route string, result interface{}) error {
	req, err := http.NewRequest("GET",fmt.Sprintf("%s/%s", c.url, route), nil)
	if err != nil{
		return fmt.Errorf("consturct http request error: %v\n", err)
	}
	res, err := c.conn.Do(req)
	if err != nil {
		return fmt.Errorf("connected error: %v\n", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response Response
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return fmt.Errorf("json unmarshal resbody error: %v\n", err)
	}
	err = json.Unmarshal(response.Result, result)
	if err != nil {
		return fmt.Errorf("json Unmarshal response result error: %v\n", err)
	}
	return nil
}

func (c *EthConn) post(route string, result interface{}, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil{
		fmt.Errorf("json marshal error: %v", err)
	}
	req, err := http.NewRequest("POST",fmt.Sprintf("%s/%s", c.url, route), bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")
	if err != nil{
		return fmt.Errorf("consturct http request error: %v\n", err)
	}
	res, err := c.conn.Do(req)
	if err != nil {
		return fmt.Errorf("connected error: %v\n", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	switch res.StatusCode {
	case 500:
		return fmt.Errorf("Server internal occured error: %s", string(resBody))
	case 400:
		return fmt.Errorf("Bad Request error: %s", string(resBody))
	}
	if err != nil {
		return err
	}
	var response Response
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return fmt.Errorf("json unmarshal resbody error: %v\n", err)
	}
	err = json.Unmarshal(response.Result, result)
	if err != nil {
		return fmt.Errorf("json Unmarshal response result error: %v\n", err)
	}
	return nil
}

func (c *EthConn) GetBlockNumber() (blockNum *big.Int, err error){
	blockNum = big.NewInt(0)
	var resStr string
	err = c.get("block", &resStr)
	blockNum.SetString(resStr, 10)
	return
}


func (c *EthConn) GetBalance(addr common.Address) (balance *big.Int, err error) {
	balance = big.NewInt(0)
	var resStr string
	err = c.get(fmt.Sprintf("balance?address=%s",addr.String()), &resStr)
	balance.SetString(resStr, 10)
	return
}

func (c *EthConn) GetErc20Balance(addr common.Address) (listBalance map[string]*big.Int, err error){
	err = c.get(fmt.Sprintf("erc20balance?address=%s",addr.String()), &listBalance)
	return
}

func (c *EthConn) GetGasPrice() (gasPrice *big.Int, err error){
	gasPrice = big.NewInt(0)
	var resStr string
	err = c.get("gasprice", &resStr)
	gasPrice.SetString(resStr, 10)
	return
}

func (c *EthConn) GetNonce(addr common.Address) (nonce uint64, err error){
	var resStr string
	err = c.get(fmt.Sprintf("nonce?address=%s",addr.String()), &resStr)
	nonce = utils.HexStrToUInt64(resStr)
	return
}

func (c *EthConn) GetTransaction(txid string) (tx types.NodeTransaction, err error){
	err = c.get(fmt.Sprintf("tx?txid=%s",txid),&tx)
	return
}

func (c *EthConn) GetNormalTransactions(address common.Address) (txs []types.EsNormalTransaction, err error){
	err = c.get(fmt.Sprintf("txs?address=%s", address.String()), &txs)
	return
}

func (c *EthConn) GetInternalTransactions(address common.Address) (txs []types.EsInternalTansaction, err error){
	err = c.get(fmt.Sprintf("intxs?address=%s",address.String()), &txs)
	return
}

func (c *EthConn) GetTokenTransactions(address common.Address) (txs []types.EsErc20TokenTransaction, err error){
	err = c.get(fmt.Sprintf("tokentxs?address=%s",address.String()), &txs)
	return
}

func (c *EthConn) GetEstimateGas(tx types.TransactionRequest) (gas uint64, err error){
	var resStr string
	err = c.post("estimategas", &resStr, tx)
	gas = utils.HexStrToUInt64(resStr)
	return
}

func (c *EthConn) SendRawTransaction(data string) (txid string, err error){
	var raw types.Raw
	raw.Hex = data
	err = c.post("send", &txid, raw)
	return
}