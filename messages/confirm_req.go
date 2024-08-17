package messages

import (
	"encoding/binary"
	"io"
	"log"
)

type ConfirmReq struct{}

type HashPair struct {
	First  [32]byte
	Second [32]byte
}

func ReadConfirmReq(r io.Reader, extensions Extensions) ConfirmReq {
	for i := 0; i < extensions.ItemCount(); i++ {
		pair := &HashPair{}
		if err := binary.Read(r, binary.LittleEndian, pair); err != nil {
			log.Fatalf("Error reading hash pair: %v", err)
		}
	}

	return ConfirmReq{}
}
