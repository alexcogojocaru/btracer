package trace

import (
	"encoding/hex"
)

const DEFAULT_SPAN_BYTES_SIZE = 8

var NullSpanID SID

type SID [DEFAULT_SPAN_BYTES_SIZE]byte

type Span struct {
	TraceID      TID
	SpanID       SID
	ParentSpanID SID
	Name         string
}

func (s *Span) Start() {

}

func (s *Span) End() {

}

func (sid *SID) ToString() string {
	return hex.EncodeToString(sid[:])
}
