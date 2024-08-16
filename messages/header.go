package messages

import (
	"encoding/binary"
	"node/config"
	"node/utils"
)

type Extensions uint16

type Header struct {
	Magic        byte
	Network      byte
	VersionMax   uint8
	VersionUsing uint8
	VersionMin   uint8
	MessageType  Type
	Extensions   Extensions
}

func (e Extensions) ItemCount() int {
	var isV2 = (e & 1) != 0

	if isV2 {
		left := (e & 0xf000) >> 12
		right := (e & 0x00f0) >> 4
		return int((left << 4) | right)
	} else {
		return int((e & 0xf000) >> 12)
	}
}

func NewHeader(messageType Type, extensions Extensions) *Header {
	return &Header{
		Magic:        'R',
		Network:      config.Network.Id,
		VersionMax:   20,
		VersionUsing: 20,
		VersionMin:   20,
		MessageType:  messageType,
		Extensions:   extensions,
	}
}

func (h *Header) Serialize() []byte {
	return utils.Serialize(h, binary.LittleEndian)
}
