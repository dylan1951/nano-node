package main

import (
	"node/config"
	"node/node"
)

func main() {
	config.Load()
	node.Connect()
	node.Listen()
}
