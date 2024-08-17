package node

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"node/blocks"
	"node/config"
	"node/messages"
	"sync"
	"time"
)

type Node struct {
	peers     map[netip.Addr]*Peer
	mu        sync.Mutex
	startTime time.Time
}

func NewNode() *Node {
	node := &Node{}
	node.peers = make(map[netip.Addr]*Peer)
	node.startTime = time.Now()
	telemetry := node.Telemetry()
	fmt.Printf("%+v\n", telemetry)
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
	if addr.Addr().IsUnspecified() {
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.peers[addr.Addr()]; ok {
		fmt.Printf("Already connected to %s\n", addr.String())
		return
	}
	conn, err := net.DialTimeout("tcp", addr.String(), 2*time.Second)
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

func (n *Node) Telemetry() messages.TelemetryData {
	return messages.TelemetryData{
		NodeId:            [32]byte(config.PublicKey),
		BlockCount:        0,
		CementedCount:     0,
		UncheckedCount:    0,
		AccountCount:      0,
		BandwidthCap:      0,
		PeerCount:         uint32(len(n.peers)),
		ProtocolVersion:   config.ProtocolVersionUsing,
		Uptime:            uint64(time.Since(n.startTime).Seconds()),
		GenesisBlock:      blocks.BetaGenesisBlock.Hash(),
		MajorVersion:      0,
		MinorVersion:      0,
		PatchVersion:      0,
		PrereleaseVersion: 0,
		Maker:             2,
		Timestamp:         uint64(time.Now().UnixMilli()),
		ActiveDifficulty:  0xFFFFF00000000000,
	}
}
