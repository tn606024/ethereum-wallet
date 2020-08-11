package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ian/ethwallet/crypto"
	"github.com/ian/ethwallet/utils"
	"github.com/pborman/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
)


func NewKey() (key *crypto.Key){
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return NewKeyFromECDSA(privateKeyECDSA)
}


func NewKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *crypto.Key {
	id := uuid.NewRandom()
	key := &crypto.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

func StoreKey(keyfilepath string, key *crypto.Key, auth string) (path string,  err error) {
	keyjson, err := crypto.EncryptKey(key, auth, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return "", fmt.Errorf("encryptKey error: %v\n", err)
	}
	path, _ = filepath.Abs(keyfilepath)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("this folder doesn't exist :%s", dir)
	}
	if utils.FileExists(path) {
		return "", fmt.Errorf("this path already exist file")
	}
	err = ioutil.WriteFile(path, keyjson, 0644)
	if err != nil {
		return "", fmt.Errorf("write file error: %v", err)
	}
	return path,nil
}

func GetKey(path, auth string) (*crypto.Key, error){
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	keyjson, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file error: %v, path: %s\n", err, path)
	}
	key, err := crypto.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, fmt.Errorf("decryptKey error: %v\n", err)
	}
	return key, nil
}
