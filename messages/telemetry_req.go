package messages

import "io"

type TelemetryReq struct{}

func ReadTelemetryReq(r io.Reader, extensions Extensions) TelemetryReq {
	println("received TelemetryReq")
	return TelemetryReq{}
}
