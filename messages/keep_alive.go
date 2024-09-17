package messages

import (
	"bytes"
	"io"
	"log"
	"net/netip"
)

type KeepAlive [8]netip.AddrPort

func ReadKeepAlive(r io.Reader, extensions Extensions) KeepAlive {
	var ka KeepAlive
	buf := make([]byte, 18)

	for i := range ka {
		_, _ = io.ReadFull(r, buf)
		_ = ka[i].UnmarshalBinary(buf)
	}

	return ka
}

func (ka KeepAlive) WriteTo(w io.Writer) {
	var response bytes.Buffer

	header := NewHeader(MsgKeepAlive, 0)

	response.Write(header.Serialize())

	for _, addrPort := range ka {
		if addrPort.Addr().Is4() {
			log.Fatalf("addr should not be ipv4")
		}
		b, _ := addrPort.MarshalBinary()
		_, _ = response.Write(b)
	}

	response.WriteTo(w)
}
