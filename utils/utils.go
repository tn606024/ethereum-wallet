package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
)
const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

func IsHexStr(str string) bool{
	_, err := hex.DecodeString(str)
	if err != nil {
		return false
	}
	return true
}

func BytesToHexStr(b []byte) string{
	s := hex.EncodeToString(b)
	res := PaddingHex(s)
	return res
}

func HexToAddress(hex string) common.Address{
	return common.HexToAddress(hex)
}

func HexStrToBigInt(hexStr string) (res *big.Int){
	if len(hexStr) == 0 {
		return big.NewInt(0)
	}
	if len(hexStr) >= 2  && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	hexBytes, _:= hex.DecodeString(hexStr)
	res = big.NewInt(0).SetBytes(hexBytes)
	return
}


func HexStrToBytes(hexStr string) []byte {
	if len(hexStr) == 0 {
		return []byte{}
	}
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	hexBytes, _:= hex.DecodeString(hexStr)
	return hexBytes
}

func HexStrToUInt64(hexStr string) (res uint64){
	if len(hexStr) == 0 || hexStr == "0x"{
		return 0
	}
	if  len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}
	if len(hexStr) %2 == 1 {
		hexStr = "0" + hexStr
	}
	hexBytes, _:= hex.DecodeString(hexStr)
	if len(hexBytes) < 8 {
		hexBytes  = AddPrefixZero(hexBytes,8)
	}
	res = binary.BigEndian.Uint64(hexBytes)
	return
}



func AddPrefixZero(hex []byte, length int ) (res []byte){
	var addBytes []byte
	addlength :=  length - len(hex)
	for i := 0; i< addlength; i++ {
		addBytes = append(addBytes, byte(0))
	}
	res = append(addBytes, hex...)
	return
}

func ConcatCopy(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}

func ConcatToArray(slices ...[]byte) [][]byte {
	res := make ([][]byte, len(slices))
	for i, slice := range slices {
		if bytes.Compare(slice,[]byte{0x00}) == 0 {
			res[i] = []byte{}
		}else {
			res[i] = slice
		}
	}
	return res
}

func Uint64ToBytes(i uint64) (b []byte){
	b = make([]byte,9)
	switch {
	case i < (1 << 8):
		b[0] = byte(i)
		return b[0:1]
	case i < (1 << 16):
		b[0] = byte(i >> 8)
		b[1] = byte(i)
		return b[0:2]
	case i < (1 << 24):
		b[0] = byte(i >> 16)
		b[1] = byte(i >> 8)
		b[2] = byte(i)
		return b[0:3]
	case i < (1 << 32):
		b[0] = byte(i >> 24)
		b[1] = byte(i >> 16)
		b[2] = byte(i >> 8)
		b[3] = byte(i)
		return b[0:4]
	case i < (1 << 40):
		b[0] = byte(i >> 32)
		b[1] = byte(i >> 24)
		b[2] = byte(i >> 16)
		b[3] = byte(i >> 8)
		b[4] = byte(i)
		return b[0:5]
	case i < (1 << 48):
		b[0] = byte(i >> 40)
		b[1] = byte(i >> 32)
		b[2] = byte(i >> 24)
		b[3] = byte(i >> 16)
		b[4] = byte(i >> 8)
		b[5] = byte(i)
		return b[0:6]
	case i < (1 << 56):
		b[0] = byte(i >> 48)
		b[1] = byte(i >> 40)
		b[2] = byte(i >> 32)
		b[3] = byte(i >> 24)
		b[4] = byte(i >> 16)
		b[5] = byte(i >> 8)
		b[6] = byte(i)
		return b[0:7]
	default:
		b[0] = byte(i >> 56)
		b[1] = byte(i >> 48)
		b[2] = byte(i >> 40)
		b[3] = byte(i >> 32)
		b[4] = byte(i >> 24)
		b[5] = byte(i >> 16)
		b[6] = byte(i >> 8)
		b[7] = byte(i)
		return b[0:8]
	}
}

func UInt64ToHex(i uint64) string {
	ib := Uint64ToBytes(i)
	str := hex.EncodeToString(ib)
	return PaddingHex(str)
}

func BigIntToHex(i *big.Int) string {
	if big.NewInt(0).Cmp(i) == 0{
		return "0x0"
	}
	b := i.Bytes()
	str := hex.EncodeToString(b)
	if str[0] == 0x30 {
		str = str[1:]
	}
	return PaddingHex(str)
}

func PaddingHex(str string) string{
	if str == "" {
		return ""
	}
	if len(str) >=2 &&  str[0:2] == "0x"{
		return str
	}
	return "0x"+str
}

func IntToBigEndianBytes(input interface{}) []byte{
	switch t, v := reflect.TypeOf(input), reflect.ValueOf(input); t.Kind() {
	case reflect.Uint8:
		return  []byte{v.Interface().(uint8)}
	case reflect.Uint16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, v.Interface().(uint16))
		return buf
	case reflect.Uint32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, v.Interface().(uint32))
		return buf
	case reflect.Uint64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.Interface().(uint64))
		return buf
	case reflect.Int32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v.Interface().(int32)))
		return buf
	case reflect.Int64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v.Interface().(int64)))
		return buf
	default:
		panic(errors.New("input type is not supported"))
	}
}

func FileExists(filename string) bool {
	path, _ := filepath.Abs(filename)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


