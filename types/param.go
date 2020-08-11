package types

import (
	"fmt"
	"github.com/ian/ethwallet/utils"
)

type BlockParam string

const (
	Latest   BlockParam = "latest"
	Earliest BlockParam = "earliest"
	Pending  BlockParam = "pending"
)

func isPredefinedBlockParam(blockParam string) bool {
	return blockParam == string(Earliest) || blockParam == string(Latest) || blockParam == string(Pending)
}

func NewBlockParam(input string) (BlockParam, error) {
	if isPredefinedBlockParam(input) == true {
		return BlockParam(input), nil
	}
	if utils.IsHexStr(input) == true{
		return BlockParam(input), nil
	}
	return "", fmt.Errorf("%s is not a legal block paramater", input)
}

type Network struct {
	ChainId byte
	Name    string
}

func (n *Network) UnmarshalText(text []byte) (err error) {
	input := string(text)
	network, err := NewNetwork(input)
	if err != nil{
		return err
	}
	*n = *network
	return nil
}

func (n Network) MarshalText() (text []byte, err error){
	return []byte(n.Name), nil
}

func NewNetwork(input string) (*Network, error) {
	switch input {
	case EthereumNet.Name:
		return EthereumNet, nil
	case RopstenNet.Name:
		return RopstenNet, nil
	case RinkebyNet.Name:
		return RinkebyNet, nil
	default:
		return nil, fmt.Errorf("input is not match any network")
	}
}

var EthereumNet = &Network{
	ChainId: 1,
	Name:    "mainnet",
}

var RopstenNet = &Network{
	ChainId: 3,
	Name:    "ropsten",
}

var RinkebyNet =&Network{
	ChainId: 4,
	Name:    "rinkeby",
}

type NetworkUrl struct {
	NodeUrl 		  string 	`json:"node_url"`
	EtherscanApiUrl   string	`json:"etherscan_api_url"`
}