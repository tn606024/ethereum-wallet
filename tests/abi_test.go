package tests

import (
	"bytes"
	"github.com/tn606024/ethwallet/utils"
	"testing"
)

var funcSigTest = []struct {
	funcSig string
	ans		[]byte
}{
	{
		funcSig: "name()",
		ans: []byte{
			0x06, 0xfd, 0xde, 0x03,
		},
	},
	{
		funcSig: "symbol()",
		ans: []byte{
			0x95, 0xd8, 0x9b, 0x41,
		},
	},
	{
		funcSig: "decimals()",
		ans: []byte{
			0x31, 0x3c, 0xe5, 0x67,
		},
	},
	{
		funcSig: "totalSupply()",
		ans: []byte{
			0x18, 0x16, 0x0d, 0xdd,
		},
	},
	{
		funcSig: "balanceOf(address)",
		ans: []byte{
			0x70, 0xa0, 0x82, 0x31,
		},
	},
	{
		funcSig: "transfer(address,uint256)",
		ans: []byte{
			0xa9, 0x05, 0x9c, 0xbb,
		},
	},
	{
		funcSig: "transferFrom(address,address,uint256)",
		ans: []byte{
			0x23, 0xb8, 0x72, 0xdd,
		},
	},
	{
		funcSig: "approve(address,uint256)",
		ans: []byte{
			0x09, 0x5e, 0xa7, 0xb3,
		},
	},
	{
		funcSig: "allowance(address,address)",
		ans: []byte{
			0xdd, 0x62, 0xed, 0x3e,
		},
	},
}

func TestToMethodID(t *testing.T) {
	for _, test:= range funcSigTest {
		res := utils.ToMethodID(test.funcSig)
		if bytes.Compare(test.ans, res) != 0 {
			t.Errorf("ToMethodID error, ans is %x, but we got %x", test.ans, res)
		}
	}
}
