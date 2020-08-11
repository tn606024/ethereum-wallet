package ethclient

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ian/ethwallet/types"
	"github.com/ian/ethwallet/utils"
	"math/big"
	"sync"
)



func (c *EthereumClient) GetBlockNumber() (blockNum string, err error){
	var params []interface{}
	err = c.call("eth_blockNumber", params, &blockNum)
	if err != nil {
		return "", err
	}
	return blockNum, nil
}


func (c *EthereumClient) GetBalance(address common.Address, blockParam types.BlockParam) (balance string, err error){
	params := []interface{}{
		address,
		blockParam,
	}
	err = c.call("eth_getBalance", params, &balance)
	if err != nil {
		return "", err
	}
	return balance, nil
}

func (c *EthereumClient) GetTransaction(txid string)  (transaction types.NodeTransaction, err error){
	params := []interface{}{
		txid,
	}
	err = c.call("eth_getTransactionByHash", params, &transaction)
	return
}

func (c *EthereumClient) GetEstimateGas(transaction *types.TransactionRequest) (estimateGas string, err error){
	params := []interface{}{
		transaction,
	}
	err = c.call("eth_estimateGas", params, &estimateGas)
	if err != nil {
		return "", err
	}

	return estimateGas, nil
}

func (c *EthereumClient)  GetGasPrice() (gasPrice string, err error){
	var params []interface{}
	err = c.call("eth_gasPrice", params, &gasPrice)
	if err != nil {
		return "", err
	}
	return gasPrice, nil
}

func (c *EthereumClient) GetTransactionCount(address common.Address, blockParam types.BlockParam) (nonce string, err error){
	params := []interface{}{
		address,
		blockParam,
	}
	err = c.call("eth_getTransactionCount", params, &nonce)
	if err != nil {
		return "", err
	}
	return nonce, nil
}

func (c *EthereumClient) GetCall(transaction *types.TransactionRequest, blockParam types.BlockParam) (res string, err error){
	params := []interface{}{
		transaction,
		blockParam,
	}
	err = c.call("eth_call", params, &res)
	if err != nil {
		return"", err
	}
	return res, nil
}

func (c *EthereumClient) callVariable(contract *common.Address, method []byte) (res string, err error) {
	data := utils.EncodeABI(method)
	txr := &types.TransactionRequest{
		To:   contract.String(),
		Data: bytesToData(data),
	}
	res, err = c.GetCall(txr, types.Latest)
	if err != nil {
		return "", err
	}
	return
}

func (c *EthereumClient) SendRawTransaction(raw string) (res string, err error){
	params := []interface{}{
		raw,
	}
	err = c.call("eth_sendRawTransaction", params, &res)
	if err != nil {
		return"", err
	}
	return res, nil
}

func (c *EthereumClient) GetErc20Balance(token *types.Erc20Token, address  *common.Address, listBalance map[string]*big.Int, wg *sync.WaitGroup, errs chan error){
	defer wg.Done()
	data := utils.EncodeABI(types.Erc20FunctionInterface.BalanceOf.MethodId, address.Bytes())
	txr := &types.TransactionRequest{
		To:   token.Address.String(),
		Data: bytesToData(data),
	}
	resp, err := c.GetCall(txr, types.Latest)
	if err != nil{
		errs <- err
	}
	res := utils.HexStrToBigInt(resp)
	c.mux.Lock()
	listBalance[token.Symbol] = types.WeiToToken(res, token.Decimals)
	c.mux.Unlock()
	return
}

func (c *EthereumClient) GetErc20ListBalance(list []*types.Erc20Token, address common.Address) (map[string]*big.Int, error){
	var wg sync.WaitGroup
	errs := make(chan error)
	wgDone := make(chan bool)
	listBalance := make(map[string]*big.Int)
	wg.Add(len(list))
	for _ , token := range list{
		go c.GetErc20Balance(token, &address, listBalance, &wg, errs)
	}
	wg.Wait()
	close(wgDone)
	select {
	case <-wgDone:
		close(errs)
		return listBalance, nil
		break
	case err := <-errs:
		close(errs)
		return nil ,err
	}
	return nil, nil
}



func(c *EthereumClient) GetErc20Name(contract *common.Address, token *types.Erc20Token, wg *sync.WaitGroup, errs chan error)  {
	defer wg.Done()
	ret, err := c.callVariable(contract, types.Erc20FunctionInterface.Name.MethodId)
	if err != nil{
		errs <- err
	}
	 dec, err := utils.DecodeSingle(ret, "string")
	 res := dec.(string)
	 c.mux.Lock()
	 token.Name = res
	 c.mux.Unlock()
	if err != nil{
		errs <- err
	}
	return
}

func(c *EthereumClient) GetErc20Symbol(contract *common.Address,  token *types.Erc20Token, wg *sync.WaitGroup,errs chan error)  {
	defer wg.Done()
	ret, err := c.callVariable(contract, types.Erc20FunctionInterface.Symbol.MethodId)
	if err != nil{
		errs <- err
	}
	dec, err := utils.DecodeSingle(ret, "string")
	res := dec.(string)
	c.mux.Lock()
	token.Symbol = res
	c.mux.Unlock()
	if err != nil{
		errs <- err
	}
	return
}

func(c *EthereumClient) GetErc20Decimals(contract *common.Address, token *types.Erc20Token, wg *sync.WaitGroup,errs chan error) {
	defer wg.Done()
	ret, err := c.callVariable(contract, types.Erc20FunctionInterface.Decimals.MethodId)
	if err != nil{
		errs <- err
	}
	resb := utils.HexStrToBigInt(ret)
	if err != nil{
		errs <- err
	}
	res := int(resb.Int64())
	c.mux.Lock()
	token.Decimals = res
	c.mux.Unlock()
	return
}

func (c *EthereumClient) GetErc20Info(contract *common.Address) (token *types.Erc20Token, err error){
	var wg sync.WaitGroup
	errs := make(chan error)
	wgDone := make(chan bool)
	token = &types.Erc20Token{Address: contract}
	wg.Add(3)
	go c.GetErc20Name(contract, token, &wg, errs)
	go c.GetErc20Symbol(contract, token, &wg, errs)
	go c.GetErc20Decimals(contract, token, &wg, errs)
	wg.Wait()
	close(wgDone)
	select {
	case <-wgDone:
		close(errs)
		return token, nil
		break
	case err := <-errs:
		close(errs)
		return nil ,err
	}
	return
}

func (c *EthereumClient) GetNormalTransactions(startBlock, endBlock int, desc bool, address common.Address) (transactions []types.EsNormalTransaction, err error) {
	params := map[string]interface{}{
		"startblock": startBlock,
		"endblock":   endBlock,
		"address":    address.String(),
	}
	if desc == true {
		params["sort"] = "desc"
	} else {
		params["sort"] = "asc"
	}
	err = c.callEtherscan("account", "txlist", params, &transactions)
	return
}

func (c *EthereumClient) GetInternalTransactions(startBlock, endBlock int, desc bool, address common.Address) (transactions []types.EsInternalTansaction, err error){
	params := map[string]interface{}{
		"startblock": startBlock,
		"endblock":   endBlock,
		"address":    address.String(),
	}
	if desc == true {
		params["sort"] = "desc"
	} else {
		params["sort"] = "asc"
	}
	err = c.callEtherscan("account", "txlistinternal", params, &transactions)
	return
}


func (c *EthereumClient) GetErc20TokenTransactions(startBlock, endBlock int, desc bool, address common.Address) (transactions []types.EsErc20TokenTransaction, err error){
	params := map[string]interface{}{
		"startblock": startBlock,
		"endblock":   endBlock,
		"address":    address.String(),
	}
	if desc == true {
		params["sort"] = "desc"
	} else {
		params["sort"] = "asc"
	}
	err = c.callEtherscan("account", "tokentx", params, &transactions)
	return
}

func (c *EthereumClient) GetLogs(fromBlock int, toBlock interface{}, address common.Address, topics map[string]string) (logs []types.EsLog, err error) {
	params := map[string]interface{}{
		"fromBlock": fromBlock,
		"toBlock":   toBlock,
		"address":   address.String(),
	}

	for key, value := range topics {
		params[key] = value
	}

	err = c.callEtherscan("logs", "getLogs", params, &logs)
	return
}

//func logToTransaction()

func createErc20TranasferLog(address common.Address) map[string]string{
	topics := make(map[string]string, 4)
	topics["topic0"] = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	//topics["topic1"] = createTopicAddress(address)
	//topics["topic1_2_opr"] = "or"
	topics["topic2"] = createTopicAddress(address)
	return topics
}

func createTopicAddress(address common.Address) string{
	addr := utils.AddPrefixZero(address.Bytes(),32)
	return bytesToData(addr)
}

func bytesToData(b []byte) string {
	return utils.PaddingHex(hex.EncodeToString(b))
}