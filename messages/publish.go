package messages

import (
	"io"
	"node/blocks"
)

type Publish blocks.Block

func ReadPublish(r io.Reader, extensions Extensions) Publish {
	println("received publish")
	return blocks.Read(r)
}
