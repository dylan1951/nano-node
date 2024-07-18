package messages

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"node/config"
	"node/messages/asc_pull_ack"
	"node/messages/asc_pull_req"
	"node/messages/confirm_req"
	"node/messages/keep_alive"
	"node/messages/node_id_handshake"
	"node/messages/telemetry_ack"
	"node/messages/telemetry_req"
)

type (
	TelemetryReq    = telemetry_req.TelemetryReq
	TelemetryAck    = telemetry_ack.TelemetryAck
	NodeIdHandshake = node_id_handshake.NodeIdHandshake
	AscPullReq      = asc_pull_req.AscPullReq
	AscPullAck      = asc_pull_ack.AscPullAck
	ConfirmReq      = confirm_req.ConfirmReq
	KeepAlive       = keep_alive.KeepAlive
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

	fmt.Printf("Received message header: %+v\n", header)

	switch header.MessageType {
	case MsgNodeIdHandshake:
		return node_id_handshake.Read(reader, header.Extensions)
	case MsgConfirmReq:
		return confirm_req.Read(reader, header.Extensions)
	case MsgKeepAlive:
		return keep_alive.Read(reader, header.Extensions)
	case MsgAscPullReq:
		return asc_pull_req.Read(reader, header.Extensions)
	case MsgAscPullAck:
		return asc_pull_ack.Read(reader, header.Extensions)
	case MsgTelemetryReq:
		return telemetry_req.Read(reader, header.Extensions)
	case MsgTelemetryAck:
		return telemetry_ack.Read(reader, header.Extensions)
	default:
		log.Fatalf("Unknown message type: 0x%x", header.MessageType)
	}
}
