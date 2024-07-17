package node

var pool [][32]byte

func (n *Node) Bootstrap() {
	for {
		for _, p := range n.peers {
			telemetry, ok := <-p.RequestTelemetry()

			return
		}
	}
}

func (n *Node) updateFrontiers() {

}
