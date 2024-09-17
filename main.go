package main

import (
	"log"
	"node/config"
	"node/ledger"
	"node/node"
	"node/types/uint128"
	"node/utils"
)

type AccountRecord struct {
	Balance uint128.Uint128
}

func main() {
	config.Load()
	address := utils.PubKeyToAddress(config.PublicKey, false)
	log.Printf("Node ID: %s\n", address)
	ledger.Init()
	node.Connect()
	go node.KeepAliveSender()
	go node.Bootstrap()
	node.Listen()
}
