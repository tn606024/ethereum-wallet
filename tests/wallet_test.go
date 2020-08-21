package tests

import (
	"fmt"
	"github.com/tn606024/ethwallet/types"
	"github.com/tn606024/ethwallet/utils"
	"github.com/tn606024/ethwallet/wallet"
	"testing"
)

var config types.Config
var testwallet *wallet.EthereumWallet

func init(){
	var err error
	config, err = types.ImportConfig("./config.json")
	if err != nil {
		fmt.Errorf("Import config error: %v", err)
		return
	}
	testwallet, err = wallet.ImportEthereumWallet(TestWalletAuth.auth, TestWalletAuth.path, config)
	if err != nil{
		fmt.Errorf("Import wallet error: %v", err)
		return
	}
}

//func TestCreateNewWallet(t *testing.T) {
//	_, err := wallet.CreateNewWallet(TestWalletAuth.auth, TestWalletAuth.path, config)
//	if err != nil {
//		t.Errorf("CreateNewWallet error: %s", err)
//	}
//}

func TestWallet_SignMessage(t *testing.T) {
	msg, err := testwallet.Wallet.SignMessage([]byte("hello"))
	if err != nil {
		t.Errorf("SignMessage error: %v", err)
	}
	fmt.Println(msg)
}

func TestWallet_VerifyMessage(t *testing.T) {
	ans := wallet.VerifyMessage(utils.HexToAddress("0x51bf0b41Ba5B034f158CF1233f16bA5450F9355B"),utils.HexStrToBytes("0x61c01b1a23624f176cbc42feda9c394ce0c9c8dd80b46ab4ca3d5dfb95a4e60335ec0f8c1bcc475dfc5bdafa697b10e56c329fdf136fee4ec800898be2412d4f00"),"hello")
	if ans != true {
		t.Errorf("Verify message error")
	}
}

