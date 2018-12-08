package handlers

import (
	"../../../blockchain"
	"../../../config"
	"../../../ethereum"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

var sumRes uint

func EthereumBalance(c *gin.Context) {
	response := ethereum.GetBalance(config.GetVerifier().VerifierPublicKey)
	c.JSON(http.StatusOK, gin.H{
		"balance": response,
	})
}

func PlasmaBalance(c *gin.Context) {

	st := make([]blockchain.Input, 0)

	resp, err := http.Get("http://localhost:3001/utxo/" + config.GetVerifier().VerifierPublicKey)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &st)
	if err != nil {
		log.Println(err)
	}

	//for _, tx := range st {
	//	sumRes = (tx.Slice.End - tx.Slice.Begin) * blockchain.WeiPerCoin
	//}

	c.JSON(http.StatusOK, st)
}

func PlasmaContractAddress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"address": config.GetVerifier().PlasmaContractAddress,
	})
}

func DepositHandler(c *gin.Context) {
	result := ethereum.Deposit(c.Param("sum"))
	c.JSON(http.StatusOK, gin.H{
		"txHash": result,
	})
}

func TransferHandler(c *gin.Context) {
	address := c.Param("address")
	sum := c.Param("sum")

	c.JSON(http.StatusOK, gin.H{
		"status": address + sum,
	})
}

func ExitHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

var id = 0

type Resp struct {
	LatestBlock string `json:"lastBlock"`
}

func LatestBlockHandler(c *gin.Context) {
	st := Resp{}
	resp, err := http.Get("http://localhost:3001/status")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &st)
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, string(body))
}

func VerifiersAmountHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"verifiers_amount": "2",
	})
}

func TotalBalanceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"balance": 1677721600000000000,
	})
}

func HistoryAllHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": "0",
	})
}

func HistoryTxHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"verifiers_amount": "0",
	})
}
