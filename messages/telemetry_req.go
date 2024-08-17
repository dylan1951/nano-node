package messages

import "io"

type TelemetryReq struct{}

func ReadTelemetryReq(r io.Reader, extensions Extensions) TelemetryReq {
	return TelemetryReq{}
}
