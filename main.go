package main

import (
	"log"
	"node/config"
	"node/ledger"
	"node/node"
	"node/utils"
)

func main() {
	config.Load()
	address := utils.PubKeyToAddress(config.PublicKey, false)
	log.Printf("Node ID: %s\n", address)
	go ledger.ProcessBlocks()
	node.Connect()
	go node.Bootstrap()
	node.Listen()
}
