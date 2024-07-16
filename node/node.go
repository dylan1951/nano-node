package node

import (
	"fmt"
	"log"
	"net"
	"node/config"
	"node/peer"
)

type Node struct {
	peers []*peer.Peer
}

func NewNode() *Node {
	node := &Node{}
	return node
}

func (n *Node) Connect() {
	ips, err := net.LookupIP(config.Network.Address)
	if err != nil {
		log.Fatalf("Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		address := fmt.Sprintf("%s:%d", ip.String(), config.Network.Port)
		n.connectPeer(address)
		break
	}
}

func (n *Node) connectPeer(address string) {
	fmt.Printf("Connecting to %s\n", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Println("Connected")
	n.peers = append(n.peers, peer.NewPeer(conn, true))
}

func (n *Node) Listen() {
	address := fmt.Sprintf("%s:%d", net.IPv4zero.String(), config.Network.Port)
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
		n.peers = append(n.peers, peer.NewPeer(conn, false))
	}
}
