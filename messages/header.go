package messages

import (
	"encoding/binary"
	"node/config"
	"node/utils"
)

type Type byte

const (
	MsgNodeIdHandshake Type = 0x0a
	MsgConfirmReq      Type = 0x04
	MsgKeepAlive       Type = 0x02
	MsgTelemetryReq    Type = 0x0c
	MsgTelemetryAck    Type = 0x0d
	MsgAscPullReq      Type = 0x0e
	MsgAscPullAck      Type = 0x0f
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
