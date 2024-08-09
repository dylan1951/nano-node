package messages

import (
	"encoding/binary"
	"io"
	"log"
)

type KeepAlive []*Socket

type Socket struct {
	Address [16]byte
	Port    uint16
}

func ReadKeepAlive(reader io.Reader, extensions uint16) *KeepAlive {
	println("received keep alive")
	peers := make(KeepAlive, 8)

	for i := 0; i < 8; i++ {
		socket := &Socket{}
		if err := binary.Read(reader, binary.LittleEndian, socket); err != nil {
			log.Fatalf("Error reading hash pair: %v", err)
		}
		peers = append(peers, socket)
	}

	return &peers
}
