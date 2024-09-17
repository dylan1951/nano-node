package node

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/netip"
	"node/config"
	"node/messages"
	"node/store"
	"sync"
	"time"
)

var peers = make(map[[16]byte]*Peer)
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
	if _, ok := peers[addr.Addr().As16()]; ok {
		return
	}
	conn, err := net.DialTimeout("tcp", addr.String(), 2*time.Second)
	if err != nil {
		log.Printf("Error: %v\n", err)
		peers[addr.Addr().As16()] = nil
		return
	}
	peer := NewPeer(conn)
	peers[addr.Addr().As16()] = peer
	go peer.handleMessages()
	peer.handleNodeIdHandshake(messages.NodeIdHandshake{})
	fmt.Printf("Discovered peer %s. There's now %d live peers.\n", addr.String(), countLivePeers())
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
		peers[peer.AddrPort().Addr().As16()] = peer
		go peer.handleMessages()
	}
}

func countLivePeers() (count uint32) {
	for _, value := range peers {
		if value != nil {
			count++
		}
	}
	return
}

func KeepAliveSender() {
	ticker := time.NewTicker(config.Network.KeepAlivePeriod)
	defer ticker.Stop()

	for range ticker.C {
		keepAlive := KeepAlive()

		for _, peer := range peers {
			if peer != nil {
				peer.SendKeepAlive(keepAlive)
			}
		}
	}
}

func Telemetry() messages.TelemetryData {
	return messages.TelemetryData{
		NodeId:            [32]byte(config.PublicKey),
		BlockCount:        store.CountBlocks(),
		CementedCount:     0,
		UncheckedCount:    0,
		AccountCount:      0,
		BandwidthCap:      0,
		PeerCount:         countLivePeers(),
		ProtocolVersion:   config.ProtocolVersionUsing,
		Uptime:            uint64(time.Since(startTime).Seconds()),
		GenesisBlock:      config.Network.Genesis.Hash(),
		MajorVersion:      0,
		MinorVersion:      0,
		PatchVersion:      0,
		PrereleaseVersion: 0,
		Maker:             69,
		Timestamp:         uint64(time.Now().UnixMilli()),
		ActiveDifficulty:  config.ActiveDifficulty,
	}
}

func KeepAlive() [8]netip.AddrPort {
	var result [8]netip.AddrPort
	keys := make([][16]byte, 0, len(peers))
	for k, peer := range peers {
		if peer != nil {
			keys = append(keys, k)
		}
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	for i := 0; i < 8; i++ {
		if i < len(keys) {
			result[i] = netip.AddrPortFrom(netip.AddrFrom16(keys[i]), config.Network.Port)
		} else {
			result[i] = netip.AddrPortFrom(netip.IPv6Unspecified(), 0)
		}
	}
	return result
}
