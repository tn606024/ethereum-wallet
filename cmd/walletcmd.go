package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tn606024/ethwallet/types"
	"github.com/tn606024/ethwallet/utils"
	"github.com/tn606024/ethwallet/wallet"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	defaultFileDir = "./keystore"
	defaultWalletConfigPath = "./config.json"
	addressSubcommand = &cli.Command{
		Name:		 "address",
		Usage: 		 "get keyfile's address",
		Description: "get keyfile's address",
		ArgsUsage: 	 "<keyfile>",
		Flags: []cli.Flag{
			keyfileFlag,
		},
		Action: func(c *cli.Context) error {
			config := loadConfig()
			wallet := unlockWallet(c, config)
			fmt.Println(wallet.Key.Address.String())
			return nil
		},
	}
	createSubcommand = &cli.Command{
		Name:        "create",
		Usage:       "create a keystore file ",
		Description: "create a keystore file, you can use keystore file to do wallet command",
		ArgsUsage:   "<keystorepath>",
		Flags: []cli.Flag{
			keystorepathFlag,
		},
		Action: func(c *cli.Context) error {
			config := loadConfig()
			var keystorepath string
			if c.String("keystorepath") != ""{
				keystorepath = c.String("keystorepath")
			}else {
				keystorepath = defaultFileDir
			}
			if utils.FileExists(keystorepath) != true {
				os.Mkdir(keystorepath,0755)
			}
			passPhrase := promptPassphrase(true)
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Please give keystore a name: ")
			name, _ := reader.ReadString('\n')
			path := fmt.Sprintf("%s/%s", keystorepath, name)
			wallet, err := wallet.CreateNewWallet(passPhrase, path, config)
			if err != nil {
				fmt.Errorf("Create Wallet Failed: %s", err)
				os.Exit(1)
			}
			fmt.Printf("Wallet is create at : %s\n", wallet.Path)
			return nil
		},
	}
	verifymessageSubcommand = &cli.Command{
		Name:			"verifymessage",
		Usage:			"verify a message is valid",
		Description: 	"verify a message is valid",
		ArgsUsage:		"<address> <signature> <message> <msgfile>",
		Flags: []cli.Flag{
			addressFlag,
			signatureFlag,
			messageFlag,
			messageFileFlag,
		},
		Action: func(c *cli.Context) error {
			var err error
			var address common.Address
			addressStr := c.String("address")
			signatureStr := c.String("signature")
			msg := loadStringOrFilePath(c, "message", "msgfile")
			address = utils.HexToAddress(addressStr)
			if err != nil {
				fmt.Printf("string to address occured error: %s", err)
				os.Exit(1)
			}
			sig := utils.HexStrToBytes(signatureStr)
			ans := wallet.VerifyMessage(address, sig, string(msg))
			fmt.Println(ans)
			return nil
		},
	}
	signmessageSubcommand = &cli.Command{
		Name: 			"signmessage",
		Usage: 			"sign a message ",
		Description: 	"sign a message with keyfile and output a raw string, message can be a string(message) or file format(msgfile)",
		ArgsUsage:	 	"<keyfile> <message> <msgfile>",
		Flags: []cli.Flag{
			keyfileFlag,
			messageFlag,
			messageFileFlag,
		},
		Action: func(c *cli.Context) error {
			config := loadConfig()
			wallet := unlockWallet(c, config)
			msg := loadStringOrFilePath(c, "message", "msgfile")
			sig, err := wallet.SignMessage(msg)
			if err != nil {
				fmt.Printf("SignMessage error: %s", err)
				os.Exit(1)
			}
			fmt.Println(sig)
			return nil
		},
	}
	signTransactionSubcommand = &cli.Command{
		Name:        "signtx",
		Usage:       "sign a transaction",
		Description: "sign a transaction with keyfile and out a raw string, transaction is json format in string(transaction) or file(txfile)",
		ArgsUsage:   "<keyfile> <transaction> <txjson>",
		Flags: []cli.Flag{
			keyfileFlag,
			transactionFlag,
			txJsonFlag,
		},
		Action: func(c *cli.Context) error {
			config := loadConfig()
			wallet := unlockWallet(c, config)
			var needsigntx types.Transaction
			txbytes := loadStringOrFilePath(c,"transaction", "txjson")
			err := json.Unmarshal(txbytes, &needsigntx)
			if err != nil {
				fmt.Printf("tx unmarshal error: %s", err)
				os.Exit(1)
			}
			raw, err := wallet.SignTxToRawTx(&needsigntx)
			if err != nil {
				fmt.Printf("sign tx error: %s", err)
				os.Exit(1)
			}
			fmt.Println(raw)
			return nil
		},
	}
	WalletCommand = &cli.Command{
		Name:	"wallet",
		Usage:	"Ethereum wallet commands",
		ArgsUsage: "",
		Category: "Wallet Commands",
		Description: "you can use wallet command to operate keyfile do action without internet, create keystore, sign message, sign Transaction, verify message" +
			         "if you want to operate wallet command, you must set network in config.json",
		Subcommands: []*cli.Command{
			addressSubcommand,
			createSubcommand,
			signmessageSubcommand,
			signTransactionSubcommand,
			verifymessageSubcommand,
		},
	}
)

