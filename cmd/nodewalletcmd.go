package cmd

import (
	"fmt"
	"github.com/ian/ethwallet/utils"
	"github.com/urfave/cli/v2"
	"math/big"
	"os"
	"strconv"
)

var(
	sendetherSubcommand = &cli.Command{
		Name:		 "sendether",
		Usage: 		 "send ether to other address",
		Description: "send ether to other address, you must set keyfile, to, value(wei), gasprice and gaslimit is optional, if you don't set, " +
					 "system will auto calculate suitable value.",
		ArgsUsage: 	 "<keyfile> <to> <value> <gasprice> <gaslimit>",
		Flags: []cli.Flag{
			keyfileFlag,
			toFlag,
			valueFlag,
			gaspriceFlag,
			gaslimitFlag,
		},
		Action: func(c *cli.Context) error {
			var err error
			gasprice := big.NewInt(0)
			gaslimit := uint64(0)
			config := loadConfig()
			wallet := unlockEthereumWallet(c, config)
			sto := c.String("to")
			svalue := c.String("value")
			sgasprice := c.String("gasprice")
			sgaslimit := c.String("gaslimit")
			to := utils.HexToAddress(sto)
			value, ok := new(big.Int).SetString(svalue, 10)
			if !ok {
				fmt.Printf("valie transfer to int occured error\n")
				os.Exit(1)
			}
			if sgasprice != ""{
				gasprice, ok = new(big.Int).SetString(sgasprice, 10)
				if !ok {
					fmt.Printf("gasprice transfer to int occured error\n")
					os.Exit(1)

				}
			}
			if sgaslimit != ""{
				gaslimit, err = strconv.ParseUint(sgaslimit, 10, 64)
				if err != nil {
					fmt.Printf("gaslimt transfer to int occured error: %s\n", sgaslimit)
					os.Exit(1)
				}
			}
			txid, err:= wallet.TransferEther(&to, value,[]byte{}, gasprice, gaslimit)
			if err != nil{
				fmt.Printf("transfer ether occured error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("transaction send success")
			fmt.Printf("txid: %s\n",txid)
			fmt.Printf("you can check tx on ethersacn: %s\n", ConstructEtherscanUrl(*config.Network, txid))
			return nil
		},
	}
	sendErc20Subcommand = &cli.Command{
		Name:		 "senderc20",
		Usage: 		 "send erc20token to other address",
		Description: "send erc20token to other address, you must set keyfile, symbol(you set in erc20_list.json in config.json), to, value(wei), gasprice and gaslimit is optional," +
			   		 "if you don't set, system will auto calculate suitable value.",
		ArgsUsage: 	 "<keyfile> <symbol> <to> <value> <gasprice> <gaslimit>",
		Flags: []cli.Flag{
			keyfileFlag,
			symbolFlag,
			toFlag,
			valueFlag,
			gaspriceFlag,
			gaslimitFlag,
		},
		Action: func(c *cli.Context) error {
			var err error
			config := loadConfig()
			wallet := unlockEthereumWallet(c, config)
			gasprice := big.NewInt(0)
			gaslimit := uint64(0)
			symbol := c.String("symbol")
			token, ok := UseSymbolFindErc20Token(config.Erc20List, symbol)
			if !ok {
				fmt.Printf("erc20 symbol %s can't find\n", symbol)
				os.Exit(1)
			}
			sto := c.String("to")
			to := utils.HexToAddress(sto)
			svalue := c.String("value")
			sgasprice := c.String("gasprice")
			sgaslimit := c.String("gaslimit")
			value, ok := new(big.Int).SetString(svalue, 10)
			if !ok {
				fmt.Printf("value transfer to int occured error\n")
				os.Exit(1)
			}
			if sgasprice != ""{
				gasprice, ok = new(big.Int).SetString(sgasprice, 10)
				if !ok {
					fmt.Printf("gasprice transfer to int occured error\n")
					os.Exit(1)

				}
			}
			if sgaslimit != ""{
				gaslimit, err = strconv.ParseUint(sgaslimit, 10, 64)
				if err != nil {
					fmt.Printf("gaslimt transfer to int occured error: %s\n", err)
					os.Exit(1)
				}
			}
			txid, err := wallet.TransferErc20(token, value,&to, gasprice, gaslimit)
			if err != nil {
				fmt.Printf("transferErc20 occured error: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("txid: %s\n",txid)
			fmt.Printf("you can check tx on ethersacn: %s\n", ConstructEtherscanUrl(*config.Network, txid))
			return nil
		},
	}
	NodewalletCommand = &cli.Command{
		Name:	"nodewallet",
		Usage:	"Ethereum nodewallet commands",
		ArgsUsage: "",
		Category: "Wallet Commands",
		Description: "you can use nodewallet command to operate keyfile do action with a connected node, if you want to operate nodewallet command" +
					 "you must set network, node_url, erc20_list.json in config.json  ",
		Subcommands: []*cli.Command{
			sendetherSubcommand,
			sendErc20Subcommand,
		},
	}
)



