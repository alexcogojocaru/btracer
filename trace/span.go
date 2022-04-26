package trace

import "encoding/hex"

const DEFAULT_SPAN_BYTES_SIZE = 8

type SID [DEFAULT_SPAN_BYTES_SIZE]byte

type Span struct {
	SpanID SID
	Name   string
}

func (sid *SID) ToString() string {
	return hex.EncodeToString(sid[:])
}
