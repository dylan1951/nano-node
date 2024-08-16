package node

import (
	"crypto/rand"
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"net"
	"node/config"
	"node/messages"
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
	for {
		msg := messages.Read(p.conn)
		switch v := msg.(type) {
		case messages.NodeIdHandshake:
			p.handleNodeIdHandshake(v)
		case messages.KeepAlive:
			p.handleKeepAlive(v)
		}
	}
}

func (p *Peer) handleKeepAlive(m messages.KeepAlive) {
	/*for _, peer := range m {
		p.node.connectPeer(peer)
	}*/
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
