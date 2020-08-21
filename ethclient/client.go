package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tn606024/ethwallet/types"
	"github.com/tn606024/ethwallet/utils"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type EthereumClient struct {
	conn    	 		*http.Client
	url 		 		string
	etherscanUrl		string
	etherscanApiKey		string
	network				*types.Network
	id      	 		int
	mux			 		sync.Mutex
}

func NewEthereumClient(url string, etherscanUrl string, etherscanApiKey string, network *types.Network) *EthereumClient {
	return &EthereumClient{
		conn: &http.Client{
			Timeout: 30 * time.Second,
		},
		url: url,
		etherscanUrl: etherscanUrl,
		etherscanApiKey: etherscanApiKey,
		network: network,
		id: 0,
	}
}

func (c *EthereumClient) call(method string, params []interface{}, result interface{}) (err error){
	if c.url == "" {
		fmt.Errorf("%s's node_url is not set in config.json",c.network.Name)
	}
	jsonrpc := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      c.id,
		"method":  method,
		"params":  params,
	}

	body, err :=json.Marshal(jsonrpc)
	if err != nil {
		return fmt.Errorf("json marshal jsonrpc error: %s\n", err)
	}
	req, err := http.NewRequest("POST", c.url, bytes.NewReader(body))
	if err != nil{
		return fmt.Errorf("consturct http request error: %s\n", err)
	}
	res, err := c.conn.Do(req)
	if err != nil {
		return fmt.Errorf("connected error: %s\n", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var response types.Response
	var responseError types.ResponseError
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return fmt.Errorf("json unmarshal resbody error: %s\n", err)
	}
	if response.Result == nil {
		err = json.Unmarshal(resBody, &responseError)
		return fmt.Errorf("res error msg: %s\n",responseError.Error.Message)

	}
	err = json.Unmarshal(response.Result, result)
	if err != nil {
		return fmt.Errorf("json Unmarshal response result error: %s\n", err)
	}
	return
}

func (c *EthereumClient) callEtherscan(module, action string, params map[string]interface{}, result interface{}) (err error) {
	if c.etherscanUrl == "" {
		fmt.Errorf("%s's etherscan_api_url is not set in config.json",c.network.Name)
	}
	req, err := http.NewRequest("GET", c.etherscanUrl, http.NoBody)
	q := req.URL.Query()
	q.Add("module", module)
	q.Add("action", action)
	q.Add("apikey", c.etherscanApiKey)
	for k, v := range params {
		q.Add(k, utils.ExtractValue(v))
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := c.conn.Do(req)
	if err != nil {
		return fmt.Errorf("connected error: %s\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var response types.EsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("json Unmarshal body error: %s\n", err)
	}
	if response.Status != 1 {
		switch response.Message {
		case "No transactions found":
			break
		default:
			err = fmt.Errorf("etherscan server error: %s", response.Message)
			return
		}

	}
	err = json.Unmarshal(response.Result, result)
	if err != nil {
		err = fmt.Errorf("json unmarshal result error: %s", response.Result)
		return
	}
	return
}
