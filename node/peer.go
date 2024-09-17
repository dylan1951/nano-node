package node

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"net"
	"net/netip"
	"node/blocks"
	"node/config"
	"node/messages"
	"node/types"
	"node/utils"
	"sync"
)

type Peer struct {
	Id            ed25519.PublicKey
	cookie        []byte
	conn          net.Conn
	mu            sync.Mutex
	frontiersChan chan []*messages.Frontier
	blocksChan    chan []blocks.Block
}

func NewPeer(conn net.Conn) *Peer {
	p := new(Peer)
	p.conn = conn
	p.cookie = make([]byte, 32)
	p.frontiersChan = make(chan []*messages.Frontier)
	p.blocksChan = make(chan []blocks.Block)
	p.mu.Lock()

	if _, err := rand.Read(p.cookie); err != nil {
		log.Fatalf("Error generating random bytes: %v", err)
	}

	return p
}

func (p *Peer) RequestFrontiers(start [32]byte, count uint16) chan []*messages.Frontier {
	p.mu.Lock()
	frontiersReq := messages.FrontiersRequest(start, count)
	p.conn.Write(frontiersReq)
	return p.frontiersChan
}

func (p *Peer) RequestBlocks(start [32]byte, count uint8, startType messages.HashType) chan []blocks.Block {
	p.mu.Lock()
	blocksReq := messages.BlocksRequest(start, count, startType)
	p.conn.Write(blocksReq)
	println("sent blocks request")
	return p.blocksChan
}

func (p *Peer) SendKeepAlive(keepAlive messages.KeepAlive) {
	p.mu.Lock()
	defer p.mu.Unlock()
	keepAlive.WriteTo(p.conn)
}

func (p *Peer) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("peer disconnected:", p.AddrPort().String(), r)
			peers[p.AddrPort().Addr().As16()] = nil
		}
	}()

	for {
		msg := messages.Read(p.conn)
		switch v := msg.(type) {
		case messages.NodeIdHandshake:
			p.handleNodeIdHandshake(v)
		case messages.KeepAlive:
			p.handleKeepAlive(v)
		case messages.TelemetryReq:
			p.handleTelemetryReq(v)
		case messages.TelemetryAck:
			p.handleTelemetryAck(v)
		case messages.AscPullAck:
			p.handleAscPullAck(v)
		}
	}
}

func (p *Peer) handleAscPullAck(msg messages.AscPullAck) {
	println("handling AscPullAck")
	p.blocksChan <- msg.Blocks
	// todo: fix (can't know if blocks or frontiers)
	//if len(msg.Blocks) > 0 {
	//	println("received blocks")
	//	p.blocksChan <- msg.Blocks
	//} else if len(msg.Frontiers) > 0 {
	//	p.frontiersChan <- msg.Frontiers
	//}
	p.mu.Unlock()
}

func (p *Peer) handleTelemetryAck(msg messages.TelemetryAck) {
	//fmt.Printf("%+v\n", msg.TelemetryData)
}

func (p *Peer) handleTelemetryReq(msg messages.TelemetryReq) {
	telemetryData := Telemetry()
	signature := ed25519.Sign(config.PrivateKey, utils.Serialize(telemetryData, binary.BigEndian))

	telemetryReq := messages.TelemetryAck{
		Signature:     types.Signature(signature),
		TelemetryData: telemetryData,
	}

	header := messages.NewHeader(messages.MsgTelemetryAck, 202)

	var buf bytes.Buffer

	buf.Write(header.Serialize())
	utils.Write(&buf, telemetryReq)

	_, err := p.conn.Write(buf.Bytes())
	if err != nil {
		log.Fatalf("Error writing telemetry ack: %v", err)
	}
}

func (p *Peer) handleKeepAlive(m messages.KeepAlive) {
	for _, peer := range m {
		if peer.Addr().IsGlobalUnicast() {
			connectPeer(peer)
		}
	}
}

func (p *Peer) handleNodeIdHandshake(m messages.NodeIdHandshake) {
	response := messages.NodeIdHandshake{}

	if m.NodeIdResponse != nil {
		if !ed25519.Verify(m.Account[:], p.cookie, m.Signature[:]) {
			log.Fatalf("Bad handshake response")
		}
		p.Id = m.Account[:]
	} else {
		response.NodeIdQuery = &messages.NodeIdQuery{Cookie: [32]byte(p.cookie)}
	}

	if m.NodeIdQuery != nil {
		signature := ed25519.Sign(config.PrivateKey, m.Cookie[:])
		response.NodeIdResponse = &messages.NodeIdResponse{
			Account:   [32]byte(config.PublicKey),
			Signature: [64]byte(signature),
		}
	}

	response.WriteTo(p.conn)

	if p.Id != nil {
		// handshake completed
		p.mu.Unlock()
	}
}

func (p *Peer) AddrPort() netip.AddrPort {
	return p.conn.RemoteAddr().(*net.TCPAddr).AddrPort()
}

func (p *Peer) IsInbound() bool {
	return p.conn.LocalAddr().(*net.TCPAddr).AddrPort().Port() == config.Network.Port
}
