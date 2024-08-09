package messages

import (
	"encoding/binary"
	"io"
	"log"
	"node/config"
)

type Message interface{}

func Read(reader io.Reader) Message {
	header := &Header{}
	if err := binary.Read(reader, binary.LittleEndian, header); err != nil {
		log.Fatalf("Error reading message header: %v\n", err)
	}

	if header.Magic != 'R' || header.Network != config.Network.Id {
		log.Fatalf("Invalid message header")
	}

	if header.VersionMax < 20 {
		log.Fatalf("Protocol version too low")
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
	default:
		log.Fatalf("Unknown message type: 0x%x", header.MessageType)
	}

	return nil
}
