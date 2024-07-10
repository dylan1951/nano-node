package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/accept-nano/ed25519-blake2b"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
)

var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey

var myCookie []byte

type MessageType byte

const (
	NodeIdHandshake MessageType = 0x0a
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

type NodeIdResponse struct {
	Account   [32]byte
	Signature [64]byte
}

func NewMessageHeader(messageType MessageType) *MessageHeader {
	return &MessageHeader{
		Magic:        'R',
		Network:      'X',
		VersionMax:   19,
		VersionUsing: 20,
		VersionMin:   18,
		MessageType:  messageType,
		Extensions:   1,
	}
}

func (mh *MessageHeader) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, mh); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func handshakeMessage() ([]byte, error) {
	messageHeader := NewMessageHeader(NodeIdHandshake)

	headerBytes, err := messageHeader.Serialize()
	if err != nil {
		return nil, err
	}

	myCookie = make([]byte, 32)
	if _, err := rand.Read(myCookie); err != nil {
		log.Fatalf("Error generating random bytes: %v", err)
	}

	buffer := append(headerBytes, myCookie...)
	return buffer, nil
}

func connect(address string) {
	conn, err := net.Dial("tcp", "168.119.169.220:17075")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Connected")
	defer conn.Close()

	go receiveMessages(conn)

	buffer, err := handshakeMessage()
	if err != nil {
		log.Fatalf("Error creating handshake message: %v", err)
	}

	if _, err := conn.Write(buffer); err != nil {
		log.Fatalf("Error writing to connection: %v", err)
	}

	fmt.Println("Handshake message sent")

	select {}
}

func receiveMessages(conn net.Conn) {
	for {
		header := &MessageHeader{}
		if err := binary.Read(conn, binary.LittleEndian, header); err != nil {
			log.Println("Error reading message header:", err)
			return
		}

		if header.Magic != 'R' || header.Network != 'X' {
			log.Println("Invalid message header")
			return
		}

		fmt.Printf("Received message header: %+v\n", header)

		var cookie []byte
		var idResponse *NodeIdResponse

		if header.Extensions&1 != 0 {
			cookie = make([]byte, 32)
			if _, err := conn.Read(cookie); err != nil {
				log.Println("Error reading cookie data:", err)
				return
			}
		}

		if header.Extensions&2 != 0 {
			idResponse = &NodeIdResponse{}
			if err := binary.Read(conn, binary.LittleEndian, idResponse); err != nil {
				log.Println("Error reading NodeIdResponse:", err)
				return
			}
		}

		if cookie == nil || idResponse == nil {
			log.Println("Invalid handshake data")
			return
		}

		if !ed25519.Verify(idResponse.Account[:], myCookie, idResponse.Signature[:]) {
			log.Println("Invalid signature")
			return
		}

		//signature := ed25519.Sign(privateKey, cookie)
		//response := append(signature, publicKey...)

	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	seed, err := hex.DecodeString(os.Getenv("SEED"))
	if err != nil || len(seed) != ed25519.SeedSize {
		log.Fatalf("Invalid SEED in .env file")
	}

	privateKey = ed25519.NewKeyFromSeed(seed)
	publicKey = privateKey.Public().(ed25519.PublicKey)

	fmt.Println("Hello World")

	ips, err := net.LookupIP("peering-test.nano.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		address := fmt.Sprintf("%s:17075", ip.String())
		connect(address)
		return
	}

}
