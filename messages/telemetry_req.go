package messages

import "io"

type TelemetryReq struct{}

func ReadTelemetryReq(reader io.Reader, extensions uint16) *TelemetryReq {
	println("received ReadTelemetryReq")
	return nil
}
