package cmd

import 	"github.com/urfave/cli/v2"


var (
	keystorepathFlag = &cli.StringFlag{
		Name: 	 "keystorepath",
		Aliases: []string{"p"},
		Value:	 "",
		Usage:	 "file folder to store keyfile, if you haven't target, keyfile will generate in './keystore'",
	}
	keyfileFlag = &cli.StringFlag{
		Name: 	 "keyfile",
		Value: 	 "",
		Usage:	 "file contains keystore",
	}
	addressFlag = &cli.StringFlag{
		Name:	"address",
		Value: 	 "",
		Usage: 	"ethereum address",
	}
	signatureFlag = &cli.StringFlag{
		Name:	"signature",
		Usage:	"signature",
	}
	messageFlag = &cli.StringFlag{
		Name:  "message",
		Usage: "message need to sign",
	}
	messageFileFlag = &cli.StringFlag{
		Name: 	 "msgfile",
		Usage:	 "message file",
	}
	transactionFlag = &cli.StringFlag{
		Name:	"transaction",
		Usage: 	"transaction in json format",
	}
	txJsonFlag = &cli.StringFlag{
		Name:	"txjson",
		Usage: 	"transaction json file",
	}
	rawFlag = &cli.StringFlag{
		Name:	"raw",
		Usage: 	"need send raw",
	}
	rawFileFlag = &cli.StringFlag{
		Name:	"rawfile",
		Usage:	"raw in file",
	}
	fromFlag = &cli.StringFlag{
		Name: 	"from",
		Usage:	"from address",
		Required: true,
	}
	toFlag   = &cli.StringFlag{
		Name:	"to",
		Usage:	"to address",
		Required: true,
	}
	valueFlag = &cli.StringFlag{
		Name:	"value",
		Usage:	"send value",
		Required: true,
	}
	gaspriceFlag = &cli.StringFlag{
		Name:	"gasprice",
		Usage:	"gasprice",
		Value:	 "",
	}
	gaslimitFlag = &cli.StringFlag{
		Name:	"gaslimit",
		Usage:	"gaslimit",
		Value:	 "",
	}
	symbolFlag = &cli.StringFlag{
		Name:	"symbol",
		Usage: 	"erc20 symbol",
		Required: true,
	}
	portFlag = &cli.IntFlag{
		Name:	"port",
		Usage:	"port",
		Value:	8080,
	}
	networkFlag = &cli.StringFlag{
		Name:	"network",
		Usage:  "specify ethereum network: mainet, ropsten, rinkery",
		Value:	"",
	}
)

