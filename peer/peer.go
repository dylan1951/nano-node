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
	telemetryChan chan *messages.Telemetry
	frontiersChan chan []*ascpullreq.Frontier
}

//func (p *Peer) RequestBlocks(start [32]byte, startType ascpullreq.HashType) chan blocks.Block {
//	p.mu.Lock()
//}

func (p *Peer) RequestTelemetry() chan *messages.TelemetryAck {
	p.mu.Lock()
	if p.telemetryChan != nil {
		log.Fatal("telemetry channel already set")
	}
	p.telemetryChan = make(chan *messages.Telemetry)
	msg := messages.NewHeader(messages.MsgTelemetryReq, 0).Serialize()
	if _, err := p.conn.Write(msg); err != nil {
		log.Fatalf("Error writing handshake response: %v", err)
	}
	return p.telemetryChan
}

func (p *Peer) RequestFrontiers(start [32]byte) chan []*ascpullreq.Frontier {
	p.mu.Lock()
	if p.frontiersChan != nil {
		log.Fatal("frontier channel already set")
	}
	p.frontiersChan = make(chan []*ascpullreq.Frontier)
	msg := ascpullreq.FrontiersRequest(start, 10)
	if _, err := p.conn.Write(msg); err != nil {
		log.Fatalf("Error writing frontiers request: %v", err)
	}
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
		p.handleNodeIdHandshake(0)
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
		}
	}
}

func (p *Peer) handleTelemetryAck(ack messages.TelemetryAck) {

}

func (p *Peer) handleAscPullAck(ack messages.AscPullAck) {
	p.mu.Unlock()
}

func (p *Peer) handleTelemetryReq(extensions uint16) {}

func (p *Peer) handleKeepAlive(peers messages.KeepAlive) {

}

func (p *Peer) handleConfirmReq(req messages.ConfirmReq) {

}

func (p *Peer) handleAscPullReq(req messages.AscPullReq) {

}

func (p *Peer) handleNodeIdHandshake(message messages.NodeIdHandshake) {
	if p.Id != nil {
		log.Fatalf("Handshake already completed")
	}

	if message.Cookie != nil || p.Id == nil {
		var msg []byte
		extensions := uint16(0)

		if p.Id == nil {
			msg = append(msg, p.cookie...)
			extensions |= 1
		}

		if message.Cookie != nil {
			signature := ed25519.Sign(config.PrivateKey, cookie)
			msg = append(msg, config.PublicKey...)
			msg = append(msg, signature...)
			extensions |= 2
		}

		header := messages.NewHeader(messages.MsgNodeIdHandshake, extensions).Serialize()
		msg = append(header, msg...)

		if _, err := p.conn.Write(msg); err != nil {
			log.Fatalf("Error writing handshake response: %v", err)
		}
	}

	if p.Id != nil {
		p.mu.Unlock()
	}
}
