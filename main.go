package main

import (
	"node/config"
	"node/node"
)

func main() {
	config.Load()
	n := node.NewNode()
	n.Connect()
	//go n.Bootstrap()
	n.Listen()
}
