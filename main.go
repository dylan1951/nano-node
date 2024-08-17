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
	log.Printf("representative address: %s\n", address)
	node.Connect()
	node.Listen()
}
