package peer

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"net"
	"node/blocks"
	"node/config"
	"node/message"
	"node/message/ascpullreq"
	"node/utils"
)

type Peer struct {
	Conn       net.Conn
	cookie     []byte
	id         ed25519.PublicKey
	Ready      bool
	blocksChan chan blocks.Block
}

func NewPeer(conn net.Conn, blocks chan blocks.Block, weInitiated bool) *Peer {
	p := new(Peer)
	p.Conn = conn
	p.cookie = make([]byte, 32)
	p.blocksChan = blocks

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
		header := &message.Header{}
		if err := binary.Read(p.Conn, binary.LittleEndian, header); err != nil {
			log.Fatalf("Error reading message header: %v\n", err)
		}

		if header.Magic != 'R' || header.Network != config.Network.Id {
			log.Fatalf("Invalid message header")
		}

		if header.VersionMax < 20 {
			log.Println("Protocol version too low")
			break
		}

		fmt.Printf("Received message header: %+v\n", header)

		switch header.MessageType {
		case message.NodeIdHandshake:
			p.handleNodeIdHandshake(header.Extensions)
		case message.ConfirmReq:
			p.handleConfirmReq(header.Extensions)
		case message.KeepAlive:
			p.handleKeepAlive(header.Extensions)
		case message.AscPullReq:
			p.handleAscPullReq(header.Extensions)
		case message.AscPullAck:
			p.handleAscPullAck(header.Extensions)
		case message.TelemetryReq:
			p.handleTelemetryReq(header.Extensions)
		default:
			log.Fatalf("Unknown message type: 0x%x", header.MessageType)
		}
	}
}

func (p *Peer) handleAscPullAck(extensions uint16) {
	header := utils.Read[ascpullreq.Header](p.Conn, binary.BigEndian)

	switch header.Type {
	case ascpullreq.Blocks:
		b := blocks.Read(p.Conn)
		for b != nil {
			hash := b.Hash()
			fmt.Println(hex.EncodeToString(hash[:]))
			b.Print()
			b = blocks.Read(p.Conn)
		}
	case ascpullreq.Frontiers:
		account := utils.Read[[32]byte](p.Conn, binary.BigEndian)
		hash := utils.Read[[32]byte](p.Conn, binary.BigEndian)

		for *account != [32]byte{} && *hash != [32]byte{} {
			fmt.Printf("frontier for %s is %s\n", hex.EncodeToString(account[:]), hex.EncodeToString(hash[:]))
			account = utils.Read[[32]byte](p.Conn, binary.BigEndian)
			hash = utils.Read[[32]byte](p.Conn, binary.BigEndian)
		}
	default:
		log.Fatalf("Unsupported AscPullAck PullType: %x", header.Type)
	}

}

func (p *Peer) handleTelemetryReq(extensions uint16) {}

func (p *Peer) handleKeepAlive(extensions uint16) {
	type Socket struct {
		Address [16]byte
		Port    uint16
	}

	var peers []*Socket

	for i := 0; i < 8; i++ {
		socket := &Socket{}
		peers = append(peers, socket)
		if err := binary.Read(p.Conn, binary.LittleEndian, socket); err != nil {
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
		if err := binary.Read(p.Conn, binary.LittleEndian, pair); err != nil {
			log.Fatalf("Error reading hash pair: %v", err)
		}
	}
}

func (p *Peer) handleAscPullReq(extensions uint16) {
	header := &ascpullreq.Header{}
	if err := binary.Read(p.Conn, binary.BigEndian, header); err != nil {
		log.Fatalf("Error reading AscPullReq header: %v", err)
	}

	switch header.Type {
	case ascpullreq.Blocks:
		payload := &ascpullreq.BlocksPayload{}
		if err := binary.Read(p.Conn, binary.BigEndian, payload); err != nil {
			log.Fatalf("Error reading Blocks payload: %v", err)
		}
		fmt.Printf("Start: %v\n", payload.Start)
		fmt.Printf("Count: %d\n", payload.Count)
		fmt.Printf("StartType: %v\n", payload.StartType)
	case ascpullreq.AccountInfo:
	case ascpullreq.Frontiers:
	default:
		log.Fatalf("Unknown AscPullReq type: %x", header.Type)
	}
}

func (p *Peer) handleNodeIdHandshake(extensions uint16) {
	if p.id != nil {
		log.Fatalf("Handshake already completed")
	}

	var cookie []byte

	if (extensions & 1) != 0 {
		cookie = make([]byte, 32)
		if _, err := p.Conn.Read(cookie); err != nil {
			log.Fatalf("Error reading cookie data: %v", err)
		}
	}

	if (extensions & 2) != 0 {
		type NodeIdResponse struct {
			Account   [32]byte
			Signature [64]byte
		}
		idResponse := &NodeIdResponse{}
		if err := binary.Read(p.Conn, binary.LittleEndian, idResponse); err != nil {
			log.Fatalf("Error reading NodeIdResponse: %v", err)
		}
		if !ed25519.Verify(idResponse.Account[:], p.cookie, idResponse.Signature[:]) {
			log.Fatalf("Invalid signature")
		}
		p.id = idResponse.Account[:]
	}

	if cookie != nil || p.id == nil {
		var msg []byte
		extensions = 0

		if p.id == nil {
			msg = append(msg, p.cookie...)
			extensions |= 1
		}

		if cookie != nil {
			signature := ed25519.Sign(config.PrivateKey, cookie)
			msg = append(msg, config.PublicKey...)
			msg = append(msg, signature...)
			extensions |= 2
		}

		header := message.NewHeader(message.NodeIdHandshake, extensions).Serialize()
		msg = append(header, msg...)

		if _, err := p.Conn.Write(msg); err != nil {
			log.Fatalf("Error writing handshake response: %v", err)
		}
	}

	if p.id != nil {
		p.Ready = true
	}
}
