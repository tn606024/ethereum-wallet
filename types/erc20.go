package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tn606024/ethwallet/utils"
	"math/big"
)

var Erc20FunctionInterface = struct {
	Name         Method
	Symbol       Method
	Decimals     Method
	totalSupply  Method
	BalanceOf    Method
	tranfer      Method
	transferFrom Method
	approve      Method
	allowance    Method
}{
	Name: Method{
		FunctionSignature: "name()",
		MethodId: []byte{
			0x06, 0xfd, 0xde, 0x03,
		},
	},
	Symbol: Method{
		FunctionSignature: "symbol()",
		MethodId: []byte{
			0x95, 0xd8, 0x9b, 0x41,
		},
	},
	Decimals: Method{
		FunctionSignature: "decimals()",
		MethodId: []byte{
			0x31, 0x3c, 0xe5, 0x67,
		},
	},
	totalSupply: Method{
		FunctionSignature: "totalSupply()",
		MethodId: []byte{
			0x18, 0x16, 0x0d, 0xdd,
		},
	},
	BalanceOf: Method{
		FunctionSignature: "balanceOf(address)",
		MethodId: []byte{
			0x70, 0xa0, 0x82, 0x31,
		},
	},
	tranfer: Method{
		FunctionSignature: "transfer(address,uint256)",
		MethodId: []byte{
			0xa9, 0x05, 0x9c, 0xbb,
		},
	},
	transferFrom: Method{
		FunctionSignature: "transferFrom(address,address,uint256)",
		MethodId: []byte{
			0x23, 0xb8, 0x72, 0xdd,
		},
	},
	approve: Method{
		FunctionSignature: "approve(address,uint256)",
		MethodId: []byte{
			0x09, 0x5e, 0xa7, 0xb3,
		},
	},
	allowance: Method{
		FunctionSignature: "allowance(address,address)",
		MethodId: []byte{
			0xdd, 0x62, 0xed, 0x3e,
		},
	},
}

type Method struct {
	FunctionSignature string
	MethodId []byte
}

type Erc20Token struct {
	Decimals 	int      `json:"decimals"`
	Name		string   `json:"name"`
	Symbol		string `json:"symbol"`
	Address 	*common.Address  `json:"address"`
}

func (t *Erc20Token) String() (string, error) {
	ts, err := json.MarshalIndent(t,"","	")
	if err != nil {
		return "", err
	}
	return string(ts) + "\n", nil
}

func (t *Erc20Token) GenerateTransferData(value *big.Int, from *common.Address) []byte{
	data := utils.EncodeABI(Erc20FunctionInterface.tranfer.MethodId, from.Bytes(), TokenToWei(value, t.Decimals).Bytes())
	return data
}

func TokenToWei(value *big.Int, decimals int) *big.Int{
	return value.Mul(value,big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)),nil))
}

func WeiToToken(value *big.Int, decimals int) *big.Int {
	return value.Div(value,big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)),nil))
}