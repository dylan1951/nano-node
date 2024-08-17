package messages

import (
	"io"
	"node/blocks"
)

type Publish blocks.Block

func ReadPublish(r io.Reader, extensions Extensions) Publish {
	return blocks.Read(r)
}
