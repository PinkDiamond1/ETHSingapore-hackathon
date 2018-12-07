package main

import (
	"../config"
	"../transactionManager"
	"./handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Assemble block ~ each second
	manager := &transactionManager.TransactionManager{}
	transactionManager.NewBlockPublisher(manager)
	transactionManager.NewEventMonitor(manager)

	r := gin.New()
	r.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)
	handlers.Manager = manager // todo refactor this
	r.POST("/tx", handlers.PostTransaction)
	//r.GET("/config") // returns contract address and abi
	//r.GET("/status") // returns last plasma block number etc.
	r.GET("/utxo/:address", handlers.GetUtxos) // returns a list of utxos for an address

	err := r.Run(fmt.Sprintf(":%d", config.GetOperator().OperatorPort))
	if err != nil {
		log.Fatal(err)
	}

	println("Operator started")
}
