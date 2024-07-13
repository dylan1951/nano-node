package peer

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"net"
	"node/config"
)

type Peer struct {
	config config.Config
	conn   net.Conn
	cookie []byte
	id     ed25519.PublicKey
}

type MessageType byte

const (
	NodeIdHandshake MessageType = 0x0a
	ConfirmReq      MessageType = 0x04
	KeepAlive       MessageType = 0x02
)

type MessageHeader struct {
	Magic        byte
	Network      byte
	VersionMax   uint8
	VersionUsing uint8
	VersionMin   uint8
	MessageType  MessageType
	Extensions   uint16
}

func (p *Peer) NewMessageHeader(messageType MessageType, extensions uint16) []byte {
	header := &MessageHeader{
		Magic:        'R',
		Network:      p.config.Network.Id,
		VersionMax:   19,
		VersionUsing: 20,
		VersionMin:   18,
		MessageType:  messageType,
		Extensions:   extensions,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, header); err != nil {
		log.Fatalf("binary.Write failed: %s", err)
	}

	return buf.Bytes()
}

func NewPeer(config config.Config, conn net.Conn, weInitiated bool) *Peer {
	p := new(Peer)
	p.config = config
	p.conn = conn
	p.cookie = make([]byte, 32)

	if _, err := rand.Read(p.cookie); err != nil {
		log.Fatalf("Error generating random bytes: %v", err)
	}

	go p.handleMessages()

	fmt.Println(conn.RemoteAddr())
	fmt.Println(conn.LocalAddr())

	if weInitiated {
		message := p.NewMessageHeader(NodeIdHandshake, 1)
		message = append(message, p.cookie...)
		if _, err := conn.Write(message); err != nil {
			log.Fatalf("Error writing to connection: %v", err)
		}
		fmt.Println("Handshake message sent")
	}

	return p
}

func (p *Peer) handleMessages() {
	for {
		header := &MessageHeader{}
		if err := binary.Read(p.conn, binary.LittleEndian, header); err != nil {
			log.Fatalf("Error reading message header: %v\n", err)
		}

		if header.Magic != 'R' || header.Network != p.config.Network.Id {
			log.Fatalf("Invalid message header")
		}

		fmt.Printf("Received message header: %+v\n", header)

		switch header.MessageType {
		case NodeIdHandshake:
			p.handleNodeIdHandshake(header.Extensions)
		case ConfirmReq:
			p.handleConfirmReq(header.Extensions)
		case KeepAlive:
			p.handleKeepAlive(header.Extensions)
		default:
			log.Fatalf("Unknown message type: 0x%x", header.MessageType)
		}
	}
}

func (p *Peer) handleKeepAlive(extensions uint16) {
	type Socket struct {
		Address [16]byte
		Port    uint16
	}

	var peers []*Socket

	for i := 0; i < 8; i++ {
		socket := &Socket{}
		peers = append(peers, socket)
		if err := binary.Read(p.conn, binary.LittleEndian, socket); err != nil {
			log.Fatalf("Error reading hash pair: %v", err)
		}
	}
}

func (p *Peer) handleConfirmReq(extensions uint16) {
	type HashPair struct {
		First  [32]byte
		Second [32]byte
	}

	var count uint16
	var isV2 = (extensions & 1) != 0

	if isV2 {
		left := (extensions & 0xf000) >> 12
		right := (extensions & 0x00f0) >> 4
		count = (left << 4) | right
	} else {
		count = (extensions & 0xf000) >> 12
	}

	for i := 0; i < int(count); i++ {
		pair := &HashPair{}
		if err := binary.Read(p.conn, binary.LittleEndian, pair); err != nil {
			log.Fatalf("Error reading hash pair: %v", err)
		}
	}
}

func (p *Peer) handleNodeIdHandshake(extensions uint16) {
	if p.id != nil {
		log.Fatalf("Handshake already completed")
	}

	type NodeIdResponse struct {
		Account   [32]byte
		Signature [64]byte
	}

	var cookie []byte

	if (extensions & 1) != 0 {
		cookie = make([]byte, 32)
		if _, err := p.conn.Read(cookie); err != nil {
			log.Fatalf("Error reading cookie data: %v", err)
		}
	}

	if (extensions & 2) != 0 {
		idResponse := &NodeIdResponse{}
		if err := binary.Read(p.conn, binary.LittleEndian, idResponse); err != nil {
			log.Fatalf("Error reading NodeIdResponse: %v", err)
		}
		if !ed25519.Verify(idResponse.Account[:], p.cookie, idResponse.Signature[:]) {
			log.Fatalf("Invalid signature")
		}
		p.id = idResponse.Account[:]
	}

	if cookie != nil {
		message := p.NewMessageHeader(NodeIdHandshake, 2)

		if p.id == nil {
			message = p.NewMessageHeader(NodeIdHandshake, 3)
			message = append(message, p.cookie...)
		}

		signature := ed25519.Sign(p.config.PrivateKey, cookie)
		message = append(message, p.config.PublicKey...)
		message = append(message, signature...)

		if _, err := p.conn.Write(message); err != nil {
			log.Fatalf("Error writing handshake response: %v", err)
		}
	}
}
