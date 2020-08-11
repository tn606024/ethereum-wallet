package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ian/ethwallet/ethclient"
	"github.com/ian/ethwallet/types"
	"github.com/ian/ethwallet/utils"
	"net/http"
	"os"
	"strconv"
)


func SetupServer(network *types.Network, port int) *gin.Engine{
	var networkUrl *types.NetworkUrl
	path := types.LoadConfigPath()
	config, err := types.ImportConfig(path)

	erc20list := config.Erc20List
	if err != nil {
		fmt.Printf("Import Config occured error: %s", err)
		os.Exit(1)
	}
	switch network.Name {
	case "mainnet":
		networkUrl = config.Mainnet
	case "ropsten":
		networkUrl = config.Ropsten
	case "rinkeby":
		networkUrl =config.Rinkeby
	default:
		fmt.Printf("can't find netwokk url in config\n")
		os.Exit(1)
	}
	client := ethclient.NewEthereumClient(networkUrl.NodeUrl, networkUrl.EtherscanApiUrl, config.EtherscanApiKey, config.Network)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/balance", func(c *gin.Context){
		param := c.DefaultQuery("param","latest")
		addr := utils.HexToAddress(c.Query("address"))
		blockParam, err := types.NewBlockParam(param)
		if err != nil {
			c.String(http.StatusBadRequest, "param is illegal: %s", param)
			return
		}
		balance, err := client.GetBalance(addr, blockParam)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		res := utils.HexStrToBigInt(balance)
		c.JSON(http.StatusOK, gin.H{
			"result": res.String(),
		})
	})
	r.GET("/erc20balance", func(c *gin.Context){

		addr := utils.HexToAddress(c.Query("address"))
		listbalance, err := client.GetErc20ListBalance(erc20list, addr)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": listbalance,
		})
	})
	r.GET("/tx", func(c *gin.Context){
		txid := c.Query("txid")
		tx, err := client.GetTransaction(txid)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": tx,
		})
	})
	r.GET("/block", func(c *gin.Context) {
		block, err := client.GetBlockNumber()
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		res := utils.HexStrToBigInt(block)
		c.JSON(http.StatusOK, gin.H{
			"result": res.String(),
		})
	})
	r.GET("/gasprice", func(c *gin.Context) {
		gasPrice, err := client.GetGasPrice()
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
		}
		res := utils.HexStrToBigInt(gasPrice)
		c.JSON(http.StatusOK, gin.H{
			"result": res.String(),
		})
	})
	r.GET("/nonce", func(c *gin.Context) {
		addr := utils.HexToAddress( c.Query("address"))
		blockParam, err := types.NewBlockParam(c.DefaultQuery("param","latest"))
		nonce, err := client.GetTransactionCount(addr, blockParam)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		res := utils.HexStrToUInt64(nonce)
		c.JSON(http.StatusOK, gin.H{
			"result": strconv.FormatUint(res,10),
		})
	})
	r.POST("/estimategas", func(c *gin.Context){
		var txReq types.TransactionRequest
		err := c.BindJSON(&txReq)
		if err != nil{
			c.String(http.StatusBadRequest, "request is illegal")
			return
		}
		estimateGas, err := client.GetEstimateGas(&txReq)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		res := utils.HexStrToUInt64(estimateGas)
		c.JSON(http.StatusOK, gin.H{
			"result": strconv.FormatUint(res,10),
		})
	})

	r.GET("/logs", func(c *gin.Context) {
		topics := make(map[string]string)
		sFromBlock:= c.DefaultQuery("fromblock","0")
		fromBlock, err := strconv.Atoi(sFromBlock)
		if err != nil {
			c.String(http.StatusBadRequest, "fromblock is illegal, %s", sFromBlock)
			return
		}
		toBlock := c.DefaultQuery("toblock","latest")
		addr := utils.HexToAddress( c.Query("address"))
		topic0 := c.Query("topic0")
		if topic0 != "" {
			topics["topic0"] = topic0
		}
		topic1 := c.Query("topic1")
		if topic1 != "" {
			topics["topic1"] = topic1
		}
		topic2 := c.Query("topic2")
		if topic2 != "" {
			topics["topic2"] = topic2
		}
		topic3 := c.Query("topic3")
		if topic3 != "" {
			topics["topic3"] = topic3
		}
		topic0_1_opr := c.Query("topic0_1_opr")
		if topic0_1_opr != "" {
			topics["topic0_1_opr"] = topic0_1_opr
		}
		topic1_2_opr := c.Query("topic1_2_opr")
		if topic1_2_opr != "" {
			topics["topic1_2_opr"] = topic1_2_opr
		}
		topic2_3_opr := c.Query("topic2_3_opr")
		if topic2_3_opr != "" {
			topics["topic2_3_opr"] = topic2_3_opr
		}
		topic0_2_opr := c.Query("topic0_2_opr")
		if topic0_2_opr != "" {
			topics["topic0_2_opr"] = topic0_2_opr
		}
		topic0_3_opr := c.Query("topic0_3_opr")
		if topic0_3_opr != "" {
			topics["topic0_3_opr"] = topic0_3_opr
		}
		topic1_3_opr := c.Query("topic1_3_opr")
		if topic1_3_opr != "" {
			topics["topic1_3_opr"] = topic1_3_opr
		}
		logs, err := client.GetLogs(fromBlock, toBlock, addr, topics)
		if err != nil{
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": logs,
		})
	})

	r.GET("/txs", func(c *gin.Context){
		addr := utils.HexToAddress( c.Query("address"))
		sStartBlock := c.DefaultQuery("startblock","0")
		sEndBlock := c.DefaultQuery("endblock","9999999")
		sdesc := c.DefaultQuery("desc","true")
		startBlock, err := strconv.Atoi(sStartBlock)
		if err != nil{
			c.String(http.StatusBadRequest, "startBlock is illegal, %s", sStartBlock)
			return
		}
		endBlock, err := strconv.Atoi(sEndBlock)
		if err != nil{
			c.String(http.StatusBadRequest, "endBlock is illegal, %s", sEndBlock)
			return
		}
		desc, err := strconv.ParseBool(sdesc)
		if err != nil {
			c.String(http.StatusBadRequest, "desc is illegal, %s", sdesc)
			return
		}
		transactions, err := client.GetNormalTransactions(startBlock, endBlock, desc, addr)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": transactions,
		})
	})
	r.GET("/intxs", func(c *gin.Context){
		addr := utils.HexToAddress( c.Query("address"))
		sStartBlock := c.DefaultQuery("startblock","0")
		sEndBlock := c.DefaultQuery("endblock","9999999")
		sdesc := c.DefaultQuery("desc","true")
		startBlock, err := strconv.Atoi(sStartBlock)
		if err != nil{
			c.String(http.StatusBadRequest, "startBlock is illegal, %s", sStartBlock)
			return
		}
		endBlock, err := strconv.Atoi(sEndBlock)
		if err != nil{
			c.String(http.StatusBadRequest, "endBlock is illegal, %s", sEndBlock)
			return
		}
		desc, err := strconv.ParseBool(sdesc)
		if err != nil {
			c.String(http.StatusBadRequest, "desc is illegal, %s", sdesc)
			return
		}
		transactions, err := client.GetInternalTransactions(startBlock, endBlock, desc, addr)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": transactions,
		})
	})
	r.GET("/tokentxs", func(c *gin.Context){
		addr := utils.HexToAddress( c.Query("address"))
		sStartBlock := c.DefaultQuery("startblock","0")
		sEndBlock := c.DefaultQuery("endblock","9999999")
		sdesc := c.DefaultQuery("desc","true")
		startBlock, err := strconv.Atoi(sStartBlock)
		if err != nil{
			c.String(http.StatusBadRequest, "startBlock is illegal, %s", sStartBlock)
			return
		}
		endBlock, err := strconv.Atoi(sEndBlock)
		if err != nil{
			c.String(http.StatusBadRequest, "endBlock is illegal, %s", sEndBlock)
			return
		}
		desc, err := strconv.ParseBool(sdesc)
		if err != nil {
			c.String(http.StatusBadRequest, "desc is illegal, %s", sdesc)
			return
		}
		transactions, err := client.GetErc20TokenTransactions(startBlock, endBlock, desc, addr)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": transactions,
		})
	})
	r.POST("/send", func(c *gin.Context){
		var raw types.Raw
		c.BindJSON(&raw)
		res, err := client.SendRawTransaction(raw.Hex)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": res,
		})
	})
	fmt.Printf("server run at http://127.0.0.1:%d\n", port)
	return r
}