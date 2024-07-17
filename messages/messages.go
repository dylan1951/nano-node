package messages

import (
	"encoding/binary"
	"node/config"
	"node/utils"
)

type Type byte

const (
	NodeIdHandshake Type = 0x0a
	ConfirmReq      Type = 0x04
	KeepAlive       Type = 0x02
	TelemetryReq    Type = 0x0c
	TelemetryAck    Type = 0x0d
	AscPullReq      Type = 0x0e
	AscPullAck      Type = 0x0f
)

type Header struct {
	Magic        byte
	Network      byte
	VersionMax   uint8
	VersionUsing uint8
	VersionMin   uint8
	MessageType  Type
	Extensions   uint16
}

func NewHeader(messageType Type, extensions uint16) *Header {
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
