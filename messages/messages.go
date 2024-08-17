package messages

import (
	"encoding/binary"
	"io"
	"log"
	"node/config"
)

type Type byte

const (
	MsgNodeIdHandshake Type = 0x0a
	MsgConfirmReq      Type = 0x04
	MsgConfirmAck      Type = 0x05
	MsgKeepAlive       Type = 0x02
	MsgTelemetryReq    Type = 0x0c
	MsgTelemetryAck    Type = 0x0d
	MsgAscPullReq      Type = 0x0e
	MsgAscPullAck      Type = 0x0f
	MsgPublish         Type = 0x03
)

type Message interface{}

func Read(reader io.Reader) Message {
	header := &Header{}
	if err := binary.Read(reader, binary.LittleEndian, header); err != nil {
		panic("Error reading message header")
	}

	if header.Magic != 'R' || header.Network != config.Network.Id {
		panic("Invalid message header")
	}

	if header.VersionMax < config.ProtocolVersionMin {
		panic("Protocol version too low")
	}

	//fmt.Printf("Received message header: %+v\n", header)

	switch header.MessageType {
	case MsgNodeIdHandshake:
		return ReadNodeIdHandshake(reader, header.Extensions)
	case MsgConfirmReq:
		return ReadConfirmReq(reader, header.Extensions)
	case MsgKeepAlive:
		return ReadKeepAlive(reader, header.Extensions)
	case MsgAscPullReq:
		return ReadAscPullReq(reader, header.Extensions)
	case MsgAscPullAck:
		return ReadAscPullAck(reader, header.Extensions)
	case MsgTelemetryReq:
		return ReadTelemetryReq(reader, header.Extensions)
	case MsgTelemetryAck:
		return ReadTelemetryAck(reader, header.Extensions)
	case MsgConfirmAck:
		return ReadConfirmAck(reader, header.Extensions)
	case MsgPublish:
		return ReadPublish(reader, header.Extensions)
	default:
		log.Fatalf("Unknown message type: 0x%x", header.MessageType)
	}

	return nil
}
