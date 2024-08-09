package peer

import (
	"crypto/rand"
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"net"
	"node/config"
	"node/messages"
	"sync"
)

type Peer struct {
	Id            ed25519.PublicKey
	cookie        []byte
	conn          net.Conn
	mu            sync.Mutex
	telemetryChan chan *messages.TelemetryAck
	frontiersChan chan []*messages.Frontier
}

//func (p *Peer) RequestBlocks(start [32]byte, startType ascpullreq.HashType) chan blocks.Block {
//	p.mu.Lock()
//}

func (p *Peer) RequestTelemetry() chan *messages.TelemetryAck {
	p.mu.Lock()
	if p.telemetryChan != nil {
		log.Fatal("telemetry channel already set")
	}
	p.telemetryChan = make(chan *messages.TelemetryAck)
	msg := messages.NewHeader(messages.MsgTelemetryReq, 0).Serialize()
	if _, err := p.conn.Write(msg); err != nil {
		log.Fatalf("Error writing handshake response: %v", err)
	}
	return p.telemetryChan
}

func (p *Peer) RequestFrontiers(start [32]byte) chan []*messages.Frontier {
	p.mu.Lock()
	if p.frontiersChan != nil {
		log.Fatal("frontier channel already set")
	}
	p.frontiersChan = make(chan []*messages.Frontier)
	msg := messages.FrontiersRequest(start, 1000)
	if _, err := p.conn.Write(msg); err != nil {
		log.Fatalf("Error writing frontiers request: %v", err)
	}
	println("sent frontiers request")
	return p.frontiersChan
}

func NewPeer(conn net.Conn, weInitiated bool) *Peer {
	p := new(Peer)
	p.conn = conn
	p.cookie = make([]byte, 32)
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
		message := messages.Read(p.conn)
		switch v := message.(type) {
		case messages.TelemetryAck:
			p.handleTelemetryAck(v)
		case messages.NodeIdHandshake:
			p.handleNodeIdHandshake(v)
		case messages.AscPullAck:
			p.handleAscPullAck(v)
		}
	}
}

func (p *Peer) handleTelemetryAck(ack messages.TelemetryAck) {

}

func (p *Peer) handleAscPullAck(ack messages.AscPullAck) {
	if ack.Blocks != nil {
		// todo: implement
	} else if ack.Frontiers != nil {
		p.frontiersChan <- ack.Frontiers
	}
	p.frontiersChan = nil
	p.mu.Unlock()
}

func (p *Peer) handleTelemetryReq(extensions uint16) {}

func (p *Peer) handleKeepAlive(peers messages.KeepAlive) {

}

func (p *Peer) handleConfirmReq(req messages.ConfirmReq) {

}

func (p *Peer) handleAscPullReq(req messages.AscPullReq) {}

func (p *Peer) handleNodeIdHandshake(message messages.NodeIdHandshake) {
	response := messages.NodeIdHandshake{}

	if message.NodeIdResponse != nil {
		if !ed25519.Verify(message.Account[:], p.cookie, message.Signature[:]) {
			log.Fatalf("Bad handshake response")
		}
		p.Id = message.Account[:]
	} else {
		response.NodeIdQuery = &messages.NodeIdQuery{Cookie: [32]byte(p.cookie)}
	}

	if message.NodeIdQuery != nil {
		signature := ed25519.Sign(config.PrivateKey, message.Cookie[:])
		response.NodeIdResponse = &messages.NodeIdResponse{
			Account:   [32]byte(config.PublicKey),
			Signature: [64]byte(signature),
		}
	}

	response.WriteTo(p.conn)

	if p.Id != nil {
		println("handshake complete")
		p.mu.Unlock()
	}
}
