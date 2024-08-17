package node

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"net"
	"node/config"
	"node/messages"
	"node/types"
	"node/utils"
	"sync"
)

type Message struct {
	Peer    *Peer
	Message messages.Message
}

type Peer struct {
	Id     ed25519.PublicKey
	cookie []byte
	conn   net.Conn
	mu     sync.Mutex
	node   *Node
}

//func (p *Peer) RequestTelemetry() chan *messages.TelemetryAck {
//	p.mu.Lock()
//	if p.telemetryChan != nil {
//		log.Fatal("telemetry channel already set")
//	}
//	p.telemetryChan = make(chan *messages.TelemetryAck)
//	msg := messages.NewHeader(messages.MsgTelemetryReq, 0).Serialize()
//	if _, err := p.conn.Write(msg); err != nil {
//		log.Fatalf("Error writing handshake response: %v", err)
//	}
//	return p.telemetryChan
//}
//
//func (p *Peer) RequestFrontiers(start [32]byte) chan []*messages.Frontier {
//	p.mu.Lock()
//	if p.frontiersChan != nil {
//		log.Fatal("frontier channel already set")
//	}
//	p.frontiersChan = make(chan []*messages.Frontier)
//	msg := messages.FrontiersRequest(start, 1000)
//	if _, err := p.conn.Write(msg); err != nil {
//		log.Fatalf("Error writing frontiers request: %v", err)
//	}
//	println("sent frontiers request")
//	return p.frontiersChan
//}

func NewPeer(conn net.Conn, node *Node, weInitiated bool) *Peer {
	p := new(Peer)
	p.conn = conn
	p.cookie = make([]byte, 32)
	p.node = node
	p.mu.Lock()

	if _, err := rand.Read(p.cookie); err != nil {
		log.Fatalf("Error generating random bytes: %v", err)
	}

	go p.handleMessages()

	if weInitiated {
		p.handleNodeIdHandshake(messages.NodeIdHandshake{})
	}

	return p
}

func (p *Peer) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("peer disconnected:", r)
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
		}
	}
}

//func (p *Peer) telemetryLoop() {
//	for {
//		p.mu.Lock()
//		telemetryReq := messages.NewHeader(messages.MsgTelemetryReq, 0)
//		_, err := p.conn.Write(telemetryReq.Serialize())
//		if err != nil {
//			log.Fatalf("Error writing telemetry request: %v", err)
//		}
//		println("sent telemetry request")
//		p.mu.Unlock()
//		time.Sleep(60 * time.Second)
//	}
//}

func (p *Peer) handleTelemetryAck(msg messages.TelemetryAck) {
	//fmt.Printf("%+v\n", msg.TelemetryData)
}

func (p *Peer) handleTelemetryReq(msg messages.TelemetryReq) {
	println("handing telemetry request")
	telemetryData := p.node.Telemetry()
	signature := ed25519.Sign(config.PrivateKey, utils.Serialize(telemetryData, binary.BigEndian))

	telemetryReq := messages.TelemetryAck{
		Signature:     types.Signature(signature),
		TelemetryData: p.node.Telemetry(),
	}

	header := messages.NewHeader(messages.MsgTelemetryAck, 202)

	var buf bytes.Buffer

	buf.Write(header.Serialize())
	utils.Write(&buf, telemetryReq)

	fmt.Printf("sending telemetry ack of length: %d\n", buf.Len())

	_, err := p.conn.Write(buf.Bytes())
	if err != nil {
		log.Fatalf("Error writing telemetry ack: %v", err)
	}
}

func (p *Peer) handleKeepAlive(m messages.KeepAlive) {
	for _, peer := range m {
		p.node.connectPeer(peer)
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
		println("handshake completed")
		p.mu.Unlock()
	}
}
