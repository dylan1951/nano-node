package node

import (
	"fmt"
	"log"
	"net"
	"node/config"
	"node/peer"
)

type Node struct {
	config config.Config
	peers  []*peer.Peer
}

func NewNode(config config.Config) *Node {
	return &Node{
		config: config,
	}
}

func (n *Node) Bootstrap() {
	ips, err := net.LookupIP(n.config.Network.Address)
	if err != nil {
		log.Fatalf("Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		address := fmt.Sprintf("%s:%d", ip.String(), n.config.Network.Port)
		n.Connect(address)
		return
	}
}

func (n *Node) Connect(address string) {
	fmt.Printf("Connecting to %s\n", address)
	conn, err := net.Dial("tcp", "168.119.169.134:17075")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Println("Connected")
	n.peers = append(n.peers, peer.NewPeer(n.config, conn, true))
}

func (n *Node) Listen() {
	address := fmt.Sprintf("%s:%d", net.IPv4zero.String(), n.config.Network.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error listening: %v", err.Error())
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %v", err.Error())
		}
		n.peers = append(n.peers, peer.NewPeer(n.config, conn, false))
	}
}
