package messages

import (
	"encoding/binary"
	"io"
	"node/utils"
)

type ConfirmAck struct {
	Account   [32]byte
	Signature [64]byte
	Hashes    [][32]byte
}

func ReadConfirmAck(r io.Reader, extensions Extensions) *ConfirmAck {
	confirmAck := &ConfirmAck{}

	confirmAck.Account = *utils.Read[[32]byte](r, binary.LittleEndian)
	confirmAck.Signature = *utils.Read[[64]byte](r, binary.LittleEndian)
	_ = *utils.Read[byte](r, binary.LittleEndian)

	confirmAck.Hashes = make([][32]byte, extensions.ItemCount())

	for i := 0; i < extensions.ItemCount(); i++ {
		confirmAck.Hashes[i] = *utils.Read[[32]byte](r, binary.LittleEndian)
	}

	return confirmAck
}
