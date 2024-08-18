package main

import (
	"log"
	"node/config"
	"node/node"
	"node/utils"
)

func main() {
	config.Load()
	address := utils.PubKeyToAddress(config.PublicKey, false)
	log.Printf("Node ID: %s\n", address)
	node.Connect()
	node.Bootstrap()
	node.Listen()
}
