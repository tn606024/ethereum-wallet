package cmd

import (
	"fmt"
	"github.com/tn606024/ethwallet/server"
	"github.com/urfave/cli/v2"
)

var (
	startSubCommand = &cli.Command{
	Name:		 "start",
	Usage: 		 "start api server",
	Description: "start a simple restful api server which can query ethereum address's basic info through infura and etherscan",
	ArgsUsage: 	 "<port><network>",
	Flags: []cli.Flag{
		portFlag,
		networkFlag,
	},
	Action: func(c *cli.Context) error {
		port := c.Int("port")
		config := loadConfig()
		network := getNetwork(c, config)
		server.SetupServer(network, port).Run(fmt.Sprintf(":%d",port))
		return nil
	},
	}
	ServerCommand = &cli.Command{
		Name:	"server",
		Usage:	"Ethereum server commands",
		ArgsUsage: "",
		Category: "Wallet Commands",
		Description: "you can use server command to open a simple restful api server which can query ethereum address's basic info through infura and etherscan" +
			"if you want to operate server command, you must set etherscan_api_key, node_url, network, erc20_list.json in config.json",
		Subcommands: []*cli.Command{
			startSubCommand,
		},
	}
)