package tests

import (
	"bytes"
	"github.com/tn606024/ethwallet/rlp"
	"testing"
)

func TestEncodeList(t *testing.T) {
	for _, test := range testRlpArray {
		res := rlp.EncodeList(test.arr)
		if bytes.Compare(res, test.res) != 0 {
			t.Errorf("the ans is %x, but we got %x",test.res, res)
		}
	}

}