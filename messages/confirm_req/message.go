package confirm_req

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

func Read(reader io.Reader, extensions uint16) *ConfirmReq {
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

	return &ConfirmReq{}
}
