package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ian/ethwallet/crypto"
	"github.com/ian/ethwallet/rlp"
	"github.com/ian/ethwallet/utils"
	"math/big"
)


type TransactionRequest struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
}

func (t *TransactionRequest) String() (string, error){
	ts, err := json.MarshalIndent(t,"","	")
	if err != nil {
		return "", err
	}
	return string(ts) + "\n", nil
}

type Transaction struct {
	Nonce    uint64   `json:"nonce"`
	GasPrice *big.Int `json:"gasprice"`
	GasLimit uint64   `json:"gaslimit"`
	From     *common.Address `json:"from,omitempty"`
	To       *common.Address `json:"to"`
	Value    *big.Int `json:"value"`
	Data     Data     `json:"data"`
	V        *big.Int `json:"v"`
	R        *big.Int `json:"r"`
	S        *big.Int `json:"s"`
}


func CreateTransaction(to *common.Address, value *big.Int, data []byte) *Transaction {
	return &Transaction{
		To:       to,
		Value:    value,
		Data:     data,
	}
}

func NewTransaction(nonce uint64, gasprice *big.Int, gaslimit uint64, to *common.Address, value *big.Int, data []byte, v *big.Int, r *big.Int, s *big.Int) *Transaction {
	return &Transaction{
		Nonce:    nonce,
		GasPrice: gasprice,
		GasLimit: gaslimit,
		To:       to,
		Value:    value,
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	}
}

func (t *Transaction) ToTransactionRequest()  *TransactionRequest {
	txr := &TransactionRequest{}
	txr.From = t.From.String()
	txr.To = t.To.String()
	if t.GasPrice == big.NewInt(0) || t.GasPrice == nil {
		txr.GasPrice = ""
	}else {
		txr.GasPrice = utils.BigIntToHex(t.GasPrice)
	}
	if t.Value == nil {
		txr.Value = "0"
	}
	txr.Value = utils.BigIntToHex(t.Value)
	if t.GasLimit == 0 {
		txr.Gas = ""
	}else {
		txr.Gas = utils.UInt64ToHex(t.GasLimit)
	}
	txr.Data = utils.BytesToHexStr(t.Data)
	return txr
}

func (t *Transaction) ToByteArray() (res [][]byte) {
	return utils.ConcatToArray(
		utils.Uint64ToBytes(t.Nonce),
		t.GasPrice.Bytes(),
		utils.Uint64ToBytes(t.GasLimit),
		t.To.Bytes(),
		t.Value.Bytes(),
		t.Data,
		t.V.Bytes(),
		t.R.Bytes(),
		t.S.Bytes(),
	)
}

func (t *Transaction) ToSignHash(network *Network) (res []byte) {
	tx := utils.ConcatToArray(
		utils.Uint64ToBytes(t.Nonce),
		t.GasPrice.Bytes(),
		utils.Uint64ToBytes(t.GasLimit),
		t.To.Bytes(),
		t.Value.Bytes(),
		t.Data,
		[]byte{network.ChainId},
		[]byte{},
		[]byte{},
	)

	msgb := rlp.EncodeList(tx)
	msghash := crypto.Keccak256(msgb)
	res = msghash
	return
}

func (t *Transaction) ToRLP() (res []byte){
	tx := t.ToByteArray()
	res = rlp.EncodeList(tx)
	return
}

func (t *Transaction) ToRawTx() (res string){
	rlpraw := t.ToRLP()
	res = utils.BytesToHexStr(rlpraw)
	return
}

func (t *Transaction) String() (string, error){
	ts, err := json.MarshalIndent(t,"","	")
	if err != nil {
		return "", err
	}
	return string(ts) + "\n", nil
}