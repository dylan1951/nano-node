package node

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"node/config"
	"node/messages"
	"sync"
	"time"
)

var peers = make(map[netip.Addr]*Peer)
var mu sync.Mutex
var startTime = time.Now()

func Connect() {
	ips, err := net.LookupIP(config.Network.Address)
	if err != nil {
		log.Fatalf("Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		addr, _ := netip.AddrFromSlice(ip)
		addrPort := netip.AddrPortFrom(addr, config.Network.Port)
		connectPeer(addrPort)
	}
}

func connectPeer(addr netip.AddrPort) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := peers[addr.Addr()]; ok {
		fmt.Printf("Already connected to %s\n", addr.String())
		return
	}
	conn, err := net.DialTimeout("tcp", addr.String(), 2*time.Second)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Connected to %s\n", addr.String())
	peer := NewPeer(conn)
	peers[peer.Addr()] = peer
	go peer.handleMessages()
	peer.handleNodeIdHandshake(messages.NodeIdHandshake{})
}

func Listen() {
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
		peer := NewPeer(conn)
		peers[peer.Addr()] = peer
		go peer.handleMessages()
	}
}

func countAlivePeers() (count uint32) {
	for _, value := range peers {
		if value != nil {
			count++
		}
	}
	return
}

func Telemetry() messages.TelemetryData {
	return messages.TelemetryData{
		NodeId:            [32]byte(config.PublicKey),
		BlockCount:        0,
		CementedCount:     0,
		UncheckedCount:    0,
		AccountCount:      0,
		BandwidthCap:      0,
		PeerCount:         countAlivePeers(),
		ProtocolVersion:   config.ProtocolVersionUsing,
		Uptime:            uint64(time.Since(startTime).Seconds()),
		GenesisBlock:      config.Network.Genesis.Hash(),
		MajorVersion:      0,
		MinorVersion:      0,
		PatchVersion:      0,
		PrereleaseVersion: 0,
		Maker:             69,
		Timestamp:         uint64(time.Now().UnixMilli()),
		ActiveDifficulty:  0xFFFFF00000000000,
	}
}
