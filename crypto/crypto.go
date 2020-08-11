package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Keccak256(data []byte) []byte{
	return crypto.Keccak256(data)
}

func Sign(digestHash []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	return crypto.Sign(digestHash, priv)
}

func SigToPub(hash []byte, sig []byte)(*ecdsa.PublicKey, error){
	return crypto.SigToPub(hash, sig)
}

func PubkeyToAddress(p ecdsa.PublicKey) common.Address{
	return crypto.PubkeyToAddress(p)
}

func S256() elliptic.Curve{
	return crypto.S256()
}

func EncryptKey(key *Key, auth string, scryptN int, scryptP int) ([]byte, error){
	k := key.Key()
	return keystore.EncryptKey(&k, auth, scryptN, scryptP)
}

func DecryptKey(keyjson []byte, auth string) (*Key, error){
	key, err := keystore.DecryptKey(keyjson, auth)
	tkey := ToKey(*key)
	return &tkey, err
}