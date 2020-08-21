package wallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tn606024/ethwallet/conn"
	"github.com/tn606024/ethwallet/crypto"
	"github.com/tn606024/ethwallet/types"
	"github.com/tn606024/ethwallet/utils"
	"math/big"
)

type Wallet struct {
	Path    string
	Key     *crypto.Key
	Network *types.Network
}


func CreateNewWallet(auth, keyfilepath string,  config types.Config) (wallet *Wallet, err error){
	key := NewKey()
	path, err := StoreKey(keyfilepath, key, auth)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Path:    path,
		Key:     key,
		Network: config.Network,
	}, nil
}

func ImportWallet(auth, path string, config types.Config) (wallet *Wallet, err error){
	key, err := GetKey(path,auth)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Path:    path,
		Key:     key,
		Network: config.Network,
	}, nil

}

func (w *Wallet) SignTx(tx *types.Transaction) error {
	signhash := tx.ToSignHash(w.Network)
	sig, err := crypto.Sign(signhash, w.Key.PrivateKey)
	if err != nil {
		return err
	}
	v,r,s := deriveSignature(sig, w.Network)
	tx.V = v
	tx.R = r
	tx.S = s
	return nil
}

func (w *Wallet) SignMessage(message []byte) (string, error) {
	sig, err := crypto.Sign(signMessageHash(message), w.Key.PrivateKey)
	if err != nil {
		return "", err
	}
	return utils.BytesToHexStr(sig), nil
}

func signMessageHash(data []byte) []byte{
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func deriveSignature(sig []byte, network *types.Network) (v, r, s *big.Int){
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 35 + network.ChainId*2})
	return
}

func VerifyMessage(address common.Address, signature []byte, message string) bool{
	recoveredPubkey, err := crypto.SigToPub(signMessageHash([]byte(message)), signature)
	if err != nil || recoveredPubkey == nil {
		fmt.Errorf("signature verification failed: %s", err)
	}
	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	success := address == recoveredAddress
	return success
}

func (w *Wallet) SignTxToRawTx(tx *types.Transaction) (string, error) {
	err := w.SignTx(tx)
	if err != nil {
		return "", fmt.Errorf("SignTxToRawTx occured error:%s \n", err)
	}
	rawTx := tx.ToRawTx()
	return rawTx, nil
}


type EthereumWallet struct {
	conn      *conn.EthConn
	Wallet    *Wallet
	erc20List []*types.Erc20Token
}

func NewEthereumWallet(auth, path string, config types.Config) (*EthereumWallet, error) {
	wallet, err := CreateNewWallet(auth, path, config)
	if err != nil {
		return nil, err
	}
	return &EthereumWallet{
		conn:      conn.NewEthConn(config.ServerUrl),
		Wallet:    wallet,
		erc20List: config.Erc20List,
	}, nil
}

func ImportEthereumWallet(auth, path string, config types.Config) (*EthereumWallet, error){
	wallet, err := ImportWallet(auth, path, config)
	if err != nil {
		return nil, err
	}
	return &EthereumWallet{
		conn:      conn.NewEthConn(config.ServerUrl),
		Wallet:    wallet,
		erc20List: config.Erc20List,
	}, nil
}

func ImportLookupEthereumWallet(address common.Address, config types.Config) *EthereumWallet {
	key := &crypto.Key{
		Address: address,
	}
	wallet :=  &Wallet{
		Path:    "",
		Key:     key,
		Network: config.Network,
	}
	return &EthereumWallet{
		conn:      conn.NewEthConn(config.ServerUrl),
		Wallet:    wallet,
		erc20List: config.Erc20List,
	}
}


func ImportEmptyEthereumWallet(config types.Config) *EthereumWallet {
	key := &crypto.Key{
	}
	wallet :=  &Wallet{
		Path:    "",
		Key:     key,
		Network: config.Network,
	}
	return &EthereumWallet{
		conn:      conn.NewEthConn(config.ServerUrl),
		Wallet:    wallet,
		erc20List: config.Erc20List,
	}
}

func (ew *EthereumWallet) GetBalance() (*big.Int, error){
	balance, err := ew.conn.GetBalance(ew.Wallet.Key.Address)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}


func (ew *EthereumWallet) GetGasPrice() (*big.Int, error){
	gasPrice, err:= ew.conn.GetGasPrice()
	if err != nil {
		return big.NewInt(0), err
	}
	return gasPrice, nil
}

func (ew *EthereumWallet) GetNonce(param types.BlockParam) (uint64, error){
	nonce, err := ew.conn.GetNonce(ew.Wallet.Key.Address)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (ew *EthereumWallet) GetGasLimit(tx *types.TransactionRequest) (uint64, error) {
	gaslimit, err := ew.conn.GetEstimateGas(*tx)
	if err != nil {
		return 0, err
	}
	return gaslimit, nil
}

func (ew *EthereumWallet) GetErc20ListBalance() (map[string]*big.Int, error){
	list, err := ew.conn.GetErc20Balance(ew.Wallet.Key.Address)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (ew *EthereumWallet) GetNormalTransactionHistory() ([]types.EsNormalTransaction, error){
	txs, err := ew.conn.GetNormalTransactions(ew.Wallet.Key.Address)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (ew *EthereumWallet) GetErc20TokenTransactionHistory() ([]types.EsErc20TokenTransaction, error){
	txs, err := ew.conn.GetTokenTransactions(ew.Wallet.Key.Address)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (ew *EthereumWallet) GetInternalTransactionHistory() ([]types.EsInternalTansaction, error){
	txs, err := ew.conn.GetInternalTransactions(ew.Wallet.Key.Address)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (ew *EthereumWallet) SendRawTransaction(raw string) (string, error){
	txid, err := ew.conn.SendRawTransaction(raw)
	if err != nil {
		return "", fmt.Errorf("SendRawTransaction occured error: %s\n",err)
	}
	return txid, nil
}


func (ew *EthereumWallet) createNormalTransaction(to *common.Address, value *big.Int, data []byte, gasPrice *big.Int, gasLimit uint64) (*types.Transaction, error){
	var tx *types.Transaction
	var err error
	if  gasPrice == nil ||gasPrice.Cmp(big.NewInt(0)) == 0 {
		gasPrice, err = ew.GetGasPrice()
		if err != nil {
			return nil, fmt.Errorf("GetGasPrice occured error:%s \n", err)
		}
	}
	nonce, err := ew.GetNonce(types.Latest)
	if err != nil {
		return nil, fmt.Errorf("GetNonce occured error:%s \n", err)
	}
	tx = &types.Transaction{
		From:	  &ew.Wallet.Key.Address,
		To:		  to,
		GasPrice: gasPrice,
		Nonce:	  nonce,
		Value:	  value,
		Data: 	  data,
	}
	if gasLimit == 0 {
		txr := tx.ToTransactionRequest()
		gasLimit, err = ew.GetGasLimit(txr)
		if err != nil {
			return nil, fmt.Errorf("GetGasLimit occured error:%v \n", err)
		}
	}
	tx.GasLimit = gasLimit
	return tx, nil
}

func (ew *EthereumWallet) createErc20Transation(token *types.Erc20Token, value *big.Int, to *common.Address,gasPrice *big.Int, gasLimit uint64) (*types.Transaction, error){
	data := token.GenerateTransferData(value, to)
	tx, err := ew.createNormalTransaction(token.Address, big.NewInt(0), data , gasPrice, gasLimit)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (ew *EthereumWallet) TransferEther(to *common.Address, value *big.Int, data []byte, gasPrice *big.Int, gasLimit uint64) (txid string, err error) {
	ether, err := ew.GetBalance()
	if err != nil {
		return "", fmt.Errorf("get balance occured error: %s\n",err)
	}
	// tx, err := ew.createNormalTransaction(to, value, []byte{}, big.NewInt(0), 0)
	tx, err := ew.createNormalTransaction(to, value, data, gasPrice, gasLimit)
	if err != nil {
		return "",  fmt.Errorf("createNormalTransaction occured error: %s\n",err)
	}
	ok := checkValueEnough(tx.Value, tx.GasPrice, tx.GasLimit, ether)
	if !ok {
		return "", fmt.Errorf("your transaction's cost is bigger then ethers you own")
	}
	txid, err = ew.signAndPublishTx(tx)
	if err != nil {
		return "", fmt.Errorf("signAndPublishTx occured error:%s \n", err)
	}
	return
}

func checkValueEnough(value *big.Int, gasPrice *big.Int, gasLimit uint64, ether *big.Int) bool{
	tvalue := big.NewInt(0).Set(value)
	tgasPrice := big.NewInt(0).Set(gasPrice)
	if tvalue.Add(tvalue, tgasPrice.Mul(tgasPrice, big.NewInt(int64(gasLimit)))).Cmp(ether) == 1 {
		return false
	}
	return true
}

func (ew *EthereumWallet) TransferErc20(token *types.Erc20Token, value *big.Int, to *common.Address, gasPrice *big.Int, gasLimit uint64) (txid string, err error){
	tx, err := ew.createErc20Transation(token, value, to, gasPrice,gasLimit)
	//tx, err := ew.createErc20Transation(token, value, to, big.NewInt(0),0)
	if err != nil {
		return "",err
	}
	txid, err = ew.signAndPublishTx(tx)
	if err != nil  {
		return "", err
	}
	return
}

func (ew *EthereumWallet) signAndPublishTx(tx *types.Transaction) (txid string, err error){
	rawTx, err := ew.Wallet.SignTxToRawTx(tx)
	if err != nil {
		return "", fmt.Errorf("SignTxToRawTx occured error:%s \n", err)
	}
	txid, err = ew.SendRawTransaction(rawTx)
	if err != nil {
		return "", fmt.Errorf("SendRawTransaction occured error:%s \n", err)
	}
	return
}

