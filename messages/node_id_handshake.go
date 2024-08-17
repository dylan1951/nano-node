package messages

import (
	"bytes"
	"encoding/binary"
	"io"
	"node/utils"
)

type NodeIdQuery struct {
	Cookie [32]byte
}

type NodeIdResponse struct {
	Account   [32]byte
	Signature [64]byte
}

type NodeIdHandshake struct {
	*NodeIdQuery
	*NodeIdResponse
}

func ReadNodeIdHandshake(r io.Reader, extensions Extensions) NodeIdHandshake {
	nodeIdHandshake := NodeIdHandshake{}

	if (extensions & 1) != 0 {
		nodeIdHandshake.NodeIdQuery = utils.Read[NodeIdQuery](r, binary.LittleEndian)
	}

	if (extensions & 2) != 0 {
		nodeIdHandshake.NodeIdResponse = utils.Read[NodeIdResponse](r, binary.LittleEndian)
	}

	return nodeIdHandshake
}

func (n *NodeIdHandshake) Extensions() Extensions {
	var extensions Extensions

	if n.NodeIdQuery != nil {
		extensions |= 1
	}
	if n.NodeIdResponse != nil {
		extensions |= 2
	}

	return extensions
}

func (n *NodeIdHandshake) WriteTo(w io.Writer) {
	if n.NodeIdQuery == nil && n.NodeIdResponse == nil {
		return
	}

	var response bytes.Buffer

	header := NewHeader(MsgNodeIdHandshake, n.Extensions())

	response.Write(header.Serialize())

	if n.NodeIdQuery != nil {
		utils.Write(&response, n.NodeIdQuery)
	}

	if n.NodeIdResponse != nil {
		utils.Write(&response, n.NodeIdResponse)
	}

	response.WriteTo(w)
}
