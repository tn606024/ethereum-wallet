package rlp

import "github.com/tn606024/ethwallet/utils"

type Offset struct {
	short  byte
	long byte
}

var strOffset = &Offset{
	short: 0x80,
	long: 0xb7,
}

var arrOffset = &Offset{
	short: 0xc0,
	long: 0xf7,
}

func EncodeList(list [][]byte) (res []byte){
	var buf [][]byte
	for _, str := range list {
		estr := encodeString(str, strOffset)
		buf = append(buf, estr)
	}
	temp := utils.ConcatCopy(buf...)

	return encodeString(temp, arrOffset)
}

func encodeString(b []byte, offset *Offset) (res []byte) {
	if len(b) == 1 && b[0] <= 0x7F {
		res = append(res, b[0])
	} else {
		head := encodeStringHeader(len(b), offset)
		res = append(head, b...)
	}
	return res
}

func encodeStringHeader(size int, offset *Offset) []byte{
	var head []byte
	buf := make([]byte,9)
	if size < 56 {
		head = append(head, offset.short + byte(size))
	} else {
		sizesize := putint(buf[1:], uint64(size))
		buf[0] = offset.long + byte(sizesize)
		head = append(head, buf[:sizesize+1]...)
	}
	return head
}

// putint writes i to the beginning of b in big endian byte
// order, using the least number of bytes needed to represent i.
func putint(b []byte, i uint64) (size int) {
	switch {
	case i < (1 << 8):
		b[0] = byte(i)
		return 1
	case i < (1 << 16):
		b[0] = byte(i >> 8)
		b[1] = byte(i)
		return 2
	case i < (1 << 24):
		b[0] = byte(i >> 16)
		b[1] = byte(i >> 8)
		b[2] = byte(i)
		return 3
	case i < (1 << 32):
		b[0] = byte(i >> 24)
		b[1] = byte(i >> 16)
		b[2] = byte(i >> 8)
		b[3] = byte(i)
		return 4
	case i < (1 << 40):
		b[0] = byte(i >> 32)
		b[1] = byte(i >> 24)
		b[2] = byte(i >> 16)
		b[3] = byte(i >> 8)
		b[4] = byte(i)
		return 5
	case i < (1 << 48):
		b[0] = byte(i >> 40)
		b[1] = byte(i >> 32)
		b[2] = byte(i >> 24)
		b[3] = byte(i >> 16)
		b[4] = byte(i >> 8)
		b[5] = byte(i)
		return 6
	case i < (1 << 56):
		b[0] = byte(i >> 48)
		b[1] = byte(i >> 40)
		b[2] = byte(i >> 32)
		b[3] = byte(i >> 24)
		b[4] = byte(i >> 16)
		b[5] = byte(i >> 8)
		b[6] = byte(i)
		return 7
	default:
		b[0] = byte(i >> 56)
		b[1] = byte(i >> 48)
		b[2] = byte(i >> 40)
		b[3] = byte(i >> 32)
		b[4] = byte(i >> 24)
		b[5] = byte(i >> 16)
		b[6] = byte(i >> 8)
		b[7] = byte(i)
		return 8
	}
}