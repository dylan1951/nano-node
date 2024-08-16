package node

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"node/config"
	"sync"
)

type Node struct {
	peers map[netip.Addr]*Peer
	mu    sync.Mutex
}

func NewNode() *Node {
	node := &Node{}
	node.peers = make(map[netip.Addr]*Peer)
	return node
}

func (n *Node) Connect() {
	ips, err := net.LookupIP(config.Network.Address)
	if err != nil {
		log.Fatalf("Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		addr, _ := netip.AddrFromSlice(ip)
		addrPort := netip.AddrPortFrom(addr, config.Network.Port)
		n.connectPeer(addrPort)
	}
}

func (n *Node) connectPeer(addr netip.AddrPort) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if addr.Addr().IsUnspecified() {
		return
	}
	if _, ok := n.peers[addr.Addr()]; ok {
		fmt.Printf("Already connected to %s\n", addr.String())
		return
	}
	fmt.Printf("Connecting to %s\n", addr.String())
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Connected to %s\n", addr.String())
	n.peers[addr.Addr()] = NewPeer(conn, n, true)
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
		addr, _ := conn.RemoteAddr().(*net.TCPAddr)

		if err != nil {
			log.Fatalf("Error accepting: %v", err.Error())
		}
		n.peers[addr.AddrPort().Addr()] = NewPeer(conn, n, false)
	}
}
