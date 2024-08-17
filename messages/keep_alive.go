package messages

import (
	"io"
	"net/netip"
)

type KeepAlive [8]netip.AddrPort

func ReadKeepAlive(r io.Reader, extensions Extensions) KeepAlive {
	var ka KeepAlive
	buf := make([]byte, 18)

	for i := range ka {
		_, _ = io.ReadFull(r, buf)
		_ = ka[i].UnmarshalBinary(buf)
		//fmt.Printf("found peer: %s, %s\n", hex.EncodeToString(buf), ka[i].String())
	}

	return ka
}
