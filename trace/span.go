package trace

import (
	"encoding/hex"
	"time"
)

const DEFAULT_SPAN_BYTES_SIZE = 8

var NullSpanID SID

type SID [DEFAULT_SPAN_BYTES_SIZE]byte

type Span struct {
	TraceID      TID
	SpanID       SID
	ParentSpanID SID
	Channel      chan Span
	Name         string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
}

func (s *Span) Start() {
	s.StartTime = time.Now()
}

func (s *Span) End() {
	s.EndTime = time.Now()
	s.Duration = time.Duration(s.EndTime.Sub(s.StartTime).Milliseconds())
	s.Channel <- *s
}

func (sid *SID) ToString() string {
	return hex.EncodeToString(sid[:])
}
