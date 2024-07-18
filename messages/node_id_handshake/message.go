package node_id_handshake

import (
	"encoding/binary"
	"io"
	"log"
)

type NodeIdHandshake struct {
	Cookie    []byte
	Account   []byte
	Signature []byte
}

func Read(reader io.Reader, extensions uint16) *NodeIdHandshake {
	nodeIdHandshake := &NodeIdHandshake{}

	if (extensions & 1) != 0 {
		nodeIdHandshake.Cookie = make([]byte, 32)
		if _, err := reader.Read(nodeIdHandshake.Cookie); err != nil {
			log.Fatalf("Error reading cookie data: %v", err)
		}
	}

	if (extensions & 2) != 0 {
		type NodeIdResponse struct {
			Account   [32]byte
			Signature [64]byte
		}
		idResponse := &NodeIdResponse{}
		if err := binary.Read(reader, binary.LittleEndian, idResponse); err != nil {
			log.Fatalf("Error reading NodeIdResponse: %v", err)
		}
		nodeIdHandshake.Account = idResponse.Account[:]
		nodeIdHandshake.Signature = idResponse.Signature[:]
	}

	return nodeIdHandshake
}
