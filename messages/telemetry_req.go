package messages

import "io"

type TelemetryReq struct{}

func ReadTelemetryReq(reader io.Reader, extensions Extensions) *TelemetryReq {
	println("received ReadTelemetryReq")
	return nil
}
