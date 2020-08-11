package crypto

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type Key keystore.Key

func (k Key) Key() keystore.Key{
	return keystore.Key(k)
}

func ToKey(key keystore.Key) Key{
	return Key(key)
}

//type Address common.Address
//
//func (a Address) Address() common.Address{
//	return common.Address(a)
//}
//
//func ToAddress(addr common.Address) Address{
//	return Address(addr)
//}