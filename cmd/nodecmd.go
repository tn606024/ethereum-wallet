package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ian/ethwallet/types"
	"github.com/ian/ethwallet/wallet"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	balanceCmd = &cli.Command{
		Name:        "balance",
		Usage:       "get balance by an address",
		Description: "get balance(wei) by an address",
		ArgsUsage:   "<address>",
		Flags: []cli.Flag{
			addressFlag,
		},
		Action: func(c *cli.Context) error {
			wallet := getLookupEthereumWallet(c)
			balance, err := wallet.GetBalance()
			if err != nil {
				fmt.Printf("getbalance occured error: %s", err)
				os.Exit(1)
			}
			fmt.Printf("wei: %s\n", balance.String())
			fmt.Printf("ether: %s\n", weiToEther(balance))
			return nil
		},
	}
	nonceCmd = &cli.Command{
		Name:        "nonce",
		Usage:       "get nonce by an address",
		Description: "get nonce by an address",
		ArgsUsage:   "<address>",
		Flags: []cli.Flag{
			addressFlag,
		},
		Action: func(c *cli.Context) error {
			wallet := getLookupEthereumWallet(c)
			nonce, err := wallet.GetNonce(types.Latest)
			if err != nil {
				fmt.Printf("nonce occured error: %s", err)
				os.Exit(1)
			}
			fmt.Println(nonce)
			return nil
		},
	}
	erc20balanceCmd = &cli.Command{
		Name:        "erc20balance",
		Usage:       "get erc20 balance by an address",
		Description: "get erc20 balance by an address, you need to set erc20_list.json in config.json",
		ArgsUsage:   "<address>",
		Flags: []cli.Flag{
			addressFlag,
		},
		Action: func(c *cli.Context) error {
			wallet := getLookupEthereumWallet(c)
			listbalance, err := wallet.GetErc20ListBalance()
			if err != nil {
				fmt.Printf("getErc20ListBalance error: %s", err)
				os.Exit(1)
			}
			for k, v := range listbalance {
				fmt.Printf("%s: %s\n", k, v.String())
			}
			return nil
		},
	}
	txhistoryCmd = &cli.Command{
		Name:        "txhistory",
		Usage:       "get normal transaction history by address",
		Description: "get normal transaction history by address, you need to set etherscan_api_key in config.json",
		ArgsUsage:   "<address>",
		Flags: []cli.Flag{
			addressFlag,
		},
		Action: func(c *cli.Context) error {
			wallet := getLookupEthereumWallet(c)
			txs, err := wallet.GetNormalTransactionHistory()
			if err != nil {
				fmt.Printf("getNormalTransactionHistory error: %s", err)
				os.Exit(1)
			}
			for _, tx := range txs {
				txstr, err := tx.String()
				if err != nil {
					fmt.Printf("tx to String occured error: %s", err)
					os.Exit(1)
				}
				fmt.Println(txstr)
			}
			return nil
		},
	}
	internaltxhistoryCmd = &cli.Command{
		Name:        "internaltxhistory",
		Usage:       "get internal transaction history by address",
		Description: "get internal transaction history by address, you need to set etherscan_api_key in config.json",
		ArgsUsage:   "<address>",
		Flags: []cli.Flag{
			addressFlag,
		},
		Action: func(c *cli.Context) error {
			wallet := getLookupEthereumWallet(c)
			txs, err := wallet.GetInternalTransactionHistory()
			if err != nil {
				fmt.Printf("getInternalTransactionHistory error: %s", err)
				os.Exit(1)
			}
			for _, tx := range txs {
				txstr, err := tx.String()
				if err != nil {
					fmt.Printf("tx to String occured error: %s", err)
					os.Exit(1)
				}
				fmt.Println(txstr)
			}
			return nil
		},
	}
	erc20txhistoryCmd = &cli.Command{
		Name:        "erc20txhistory",
		Usage:       "get erc20 transaction history by address",
		Description: "get erc20 transaction history by address, you need to set etherscan_api_key in config.json",
		ArgsUsage:   "<address>",
		Flags: []cli.Flag{
			addressFlag,
		},
		Action: func(c *cli.Context) error {
			wallet := getLookupEthereumWallet(c)
			txs, err := wallet.GetErc20TokenTransactionHistory()
			if err != nil {
				fmt.Printf("getErc20TokenTransactionHistory error: %s", err)
				os.Exit(1)
			}
			for _, tx := range txs {
				txstr, err := tx.String()
				if err != nil {
					fmt.Printf("tx to String occured error: %s", err)
					os.Exit(1)
				}
				fmt.Println(txstr)
			}
			return nil
		},
	}
	gaspriceCmd = &cli.Command{
		Name:        "gasprice",
		Usage:       "get gasprice from node",
		Description: "get gasprice from node",
		ArgsUsage:   "",
		Flags: []cli.Flag{
		},
		Action: func(c *cli.Context) error {
			config := loadConfig()
			wallet := wallet.ImportEmptyEthereumWallet(config)
			gasprice, err := wallet.GetGasPrice()
			if err != nil {
				fmt.Printf("getGasPrice error: %s", err)
				os.Exit(1)
			}
			fmt.Printf("wei: %s\n", gasprice.String())
			return nil
		},
	}
	gaslimitCmd = &cli.Command{
		Name:        "gaslimit",
		Usage:       "get transaction's estimate gas from node",
		Description: "get transaction's estimate gas from node, you need to input transaction's json format in string or file",
		ArgsUsage:   "",
		Flags: []cli.Flag{
			transactionFlag,
			txJsonFlag,
		},
		Action: func(c *cli.Context) error {
			var tx types.Transaction
			config := loadConfig()
			wallet := wallet.ImportEmptyEthereumWallet(config)
			txbytes := loadStringOrFilePath(c,"transaction", "txjson")
			err := json.Unmarshal(txbytes, &tx)
			if err != nil {
				fmt.Printf("tx unmarshal error: %s\n", err)
				os.Exit(1)
			}
			gaslimit, err := wallet.GetGasLimit(tx.ToTransactionRequest())
			fmt.Printf("wei: %d\n", gaslimit)
			return nil
		},
	}
	sendrawtxCmd= &cli.Command{
		Name:			"sendrawtx",
		Usage:			"send raw tranasaction from node",
		Description: 	"send raw tranasaction from node",
		ArgsUsage: 		"<raw> <rawfile>",
		Flags: []cli.Flag{
			rawFlag,
			rawFileFlag,
		},
		Action: func(c *cli.Context) error {
			config := loadConfig()
			wallet := wallet.ImportEmptyEthereumWallet(config)
			rawb := loadStringOrFilePath(c,"raw","rawfile")
			res, err := wallet.SendRawTransaction(string(rawb))
			if err != nil {
				fmt.Printf("sendRawTransaction error: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(res)
			return nil
		},
	}



	NodeCommand = &cli.Command{
		Name:	"node",
		Usage:	"Ethereum node commands",
		ArgsUsage: "",
		Category: "Node Commands",
		Description: "you can use node command to connect node and get information about ethereum, if you want to connect node," +
					 "you must set node_url and network in config.json. if you want to get transaction history(normal, internal, erc20) by an address," +
			    	 "you also need set etherscan_api_key in config.json",
		Subcommands: []*cli.Command{
			balanceCmd,
			nonceCmd,
			erc20balanceCmd,
			txhistoryCmd,
			internaltxhistoryCmd,
			erc20txhistoryCmd,
			gaspriceCmd,
			gaslimitCmd,
			sendrawtxCmd,
		},
	}
)



