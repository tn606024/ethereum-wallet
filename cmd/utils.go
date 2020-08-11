package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ian/ethwallet/types"
	"github.com/ian/ethwallet/utils"
	"github.com/ian/ethwallet/wallet"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"syscall"
)

func unlockEthereumWallet(c *cli.Context, config types.Config) *wallet.EthereumWallet {
	keyfile, passPhrase := getKeyfileAndPassPhrase(c,config)
	wallet, err := wallet.ImportEthereumWallet(passPhrase, keyfile, config)
	if err != nil {
		fmt.Printf("Import wallet error: %s", err)
		os.Exit(1)
	}
	return wallet
}

func promptPassphrase(confirmation bool) string {
	fmt.Println("please input password:")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Failed to read password: %v\n", err)
	}
	passphrase := string(bytePassword)

	if confirmation {
		fmt.Println("please input password again:")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Failed to read password: %v\n", err)
		}
		confirm := string(bytePassword)
		if passphrase != confirm {
			fmt.Printf("Passwords do not match\n")
			os.Exit(1)
		}
	}
	return passphrase
}

func unlockWallet(c *cli.Context,  config types.Config) *wallet.Wallet {
	keyfile, passPhrase := getKeyfileAndPassPhrase(c,config)
	wallet, err := wallet.ImportWallet(passPhrase, keyfile, config)
	if err != nil {
		fmt.Printf("Import wallet error: %s", err)
		os.Exit(1)
	}
	return wallet
}

func getKeyfileAndPassPhrase(c *cli.Context,  config types.Config) (keyfile string, passPhrase string){
	if c.String("keyfile") == "" {
		if config.Keyfile == "" {
			fmt.Println("you need specify --keyfile or filled keyfile in config.json")
			os.Exit(1)
		}
		keyfile = config.Keyfile
		if config.Passphrase != "" {
			passPhrase = config.Passphrase
		}else {
			passPhrase = promptPassphrase(true)
		}
	}else {
		keyfile = c.String("keyfile")
		passPhrase = promptPassphrase(true)
	}
	return
}

func loadStringOrFilePath(c *cli.Context,  inputFlagName string,  inputFilePathFlagName string) []byte{
	var out []byte
	var err error
	inputString := c.String(inputFlagName)
	inputFilePath := c.String(inputFilePathFlagName)
	if inputString == "" {
		if inputFilePath == "" {
			fmt.Printf("need provide %s or %s\n", inputFlagName, inputFilePathFlagName)
			os.Exit(1)
		}
		out, err = ioutil.ReadFile(inputFilePath)
		if err != nil{
			fmt.Printf("read file occured error\n")
			os.Exit(1)
		}
	}else {
		out = []byte(inputString)
	}
	return out
}


func getLookupEthereumWallet(c *cli.Context) *wallet.EthereumWallet {
	config := loadConfig()
	address := getAddress(c, config)
	wallet := wallet.ImportLookupEthereumWallet(address, config)
	return wallet
}

func getAddress(c *cli.Context, config types.Config) (address common.Address){
	var saddr string
	if c.String("address") == ""{
		if config.Address == "" {
			fmt.Println("you need to specify --address or filled address in config.json")
			os.Exit(1)
		}else{
			saddr = config.Address
		}
	}else {
		saddr = c.String("address")
	}
	address = utils.HexToAddress(saddr)
	return
}

func getNetwork(c *cli.Context, config types.Config) (network *types.Network){
	var snet string
	var err error
	if c.String("network") == ""{
		if config.Network == nil {
			fmt.Println("you need to specify --network or filled network in config.json")
			os.Exit(1)
		}else{
			network = config.Network
		}
	}else {
		snet = c.String("network")
		network, err = types.NewNetwork(snet)
		if err != nil{
			fmt.Printf("network is not vaild:%s\n", network.Name)
			os.Exit(1)		}
		}
	return
}

func loadConfig() types.Config {
	path := loadConfigPath()
	config, err := types.ImportConfig(path)
	if err != nil {
		panic(err)
	}
	return config
}

func loadConfigPath() string{
	path := os.Getenv("ETHEREUM_WALLET_CONFIG_PATH")
	if path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		path = currentPath + "/config.json"
		if utils.FileExists(path) == false {
			fmt.Printf("doesn't set config\n")
			os.Exit(1)
		}
	}
	return path
}

func UseSymbolFindErc20Token(tokens []*types.Erc20Token, symbol string) (*types.Erc20Token, bool){
	for _, token := range tokens {
		if token.Symbol == symbol {
			return token, true
		}
	}
	return nil, false
}
func ConstructEtherscanUrl(network types.Network, txid string) string{
	if network.Name == types.EthereumNet.Name {
		return fmt.Sprintf("https://etherscan.io/tx/%s\n",txid)
	}else {
		return fmt.Sprintf("https://%s.etherscan.io/tx/%s\n",network.Name,txid)
	}
}

func weiToEther(wei *big.Int) string{
	fwei := new(big.Float)
	fwei.SetString(wei.String())
	ethValue := new(big.Float).Quo(fwei, big.NewFloat(math.Pow10(18)))
	return ethValue.String()
}