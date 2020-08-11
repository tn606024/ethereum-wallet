package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ian/ethwallet/utils"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Raw struct {
	Hex string
}

type BigIntHex big.Int

func (b *BigIntHex) UnmarshalText(text []byte) (err error) {
	res := utils.HexStrToBigInt(string(text))
	*b = BigIntHex(*res)
	return
}

func (b BigIntHex) MarshalText() (text []byte, err error){
	//return []byte(utils.BigIntToHex((*big.Int)(&b))), nil
	return []byte((*big.Int)(&b).String()), nil
}

type IntHex	int

func (i *IntHex) UnmarshalText(text []byte) (err error) {
	res := int(int64(utils.HexStrToUInt64(string(text))))
	*i = IntHex(res)
	return
}
func (i IntHex) MarshalText() (text []byte, err error){
	//return []byte(utils.UInt64ToHex(uint64(int64(i)))), nil
	return []byte(strconv.Itoa(int(i))), nil
}

// BigInt is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type BigInt big.Int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *BigInt) UnmarshalText(text []byte) (err error) {
	var bigInt = new(big.Int)
	err = bigInt.UnmarshalText(text)
	if err != nil {
		return
	}

	*b = BigInt(*bigInt)
	return nil
}

// MarshalText implements the encoding.TextMarshaler
func (b BigInt) MarshalText() (text []byte, err error) {
	return []byte((*big.Int)(&b).String()), nil
}


// Time is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type Time time.Time

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(text []byte) (err error) {
	input, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		err = fmt.Errorf("strconv.ParseInt %s", err)
		return
	}

	var timestamp = time.Unix(input, 0)
	*t = Time(timestamp)

	return nil
}

// Time returns t's time.Time form
func (t Time) Time() time.Time {
	return time.Time(t)
}

// MarshalText implements the encoding.TextMarshaler
func (t Time) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(t.Time().Unix(), 10)), nil
}

type Data []byte

func (d *Data) MarshalJSON() ([]byte, error) {
	res, err:= json.Marshal(hex.EncodeToString(*d))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Data) UnmarshalText(text []byte) (err error){
	out := utils.HexStrToBytes(string(text))
	*d = out
	return nil
}

type Config struct {
	Network   		*Network    	 `json:"network"`
	Rinkeby			*NetworkUrl		 `json:"rinkeby"`
	Ropsten			*NetworkUrl		 `json:"ropsten"`
	Mainnet			*NetworkUrl		 `json:"mainnet"`
	ServerUrl		string		 	 `json:"server_url"`
	Keyfile			string			 `json:"keyfile"`
	Passphrase		string			 `json:"passphrase"`
	Address			string			 `json:"address"`
	EtherscanApiKey	string      	 `json:"etherscan_api_Key"`
	Erc20List 		[]*Erc20Token 	 `json:"erc20_list"`
}

func LoadConfigPath() string{
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

func ImportConfig(path string)(Config, error){
	var config Config
	path, err := filepath.Abs(path)
	if err != nil {
		return Config{}, err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}



