package node

import (
	"log"
	"node/message/ascpullreq"
	"node/types"
)

func (n *Node) Bootstrap() {
	//start := store.GetLastBlockHash()

	for {
		for _, p := range n.peers {
			if p.Ready {
				msg := ascpullreq.BlocksRequest(types.MustParseAddress("nano_1aa4a83yg5ewn3oyr4jx3sozg5xh4bp4ttmxpu1g4j9uzzmnxckk64dcepyx"), ascpullreq.Account)

				if _, err := p.Conn.Write(msg); err != nil {
					log.Fatalf("Error writing FrontiersRequest: %v", err)
				}

				return
			}
		}
	}
}
