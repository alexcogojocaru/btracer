package trace

import (
	"encoding/hex"
	"time"

	agent "github.com/alexcogojocaru/btracer/proto-gen/btrace_proxy"
)

const DEFAULT_SPAN_BYTES_SIZE = 8

var NullSpanID SID

type SID [DEFAULT_SPAN_BYTES_SIZE]byte

type Span struct {
	TraceID      TID
	SpanID       SID
	ParentSpanID SID
	ServiceName  string
	Channel      chan Span
	Name         string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Logs         []*agent.KeyValue
}

func (s *Span) Start() {
	s.StartTime = time.Now()
}

func (s *Span) End() {
	s.EndTime = time.Now()
	s.Duration = time.Duration(s.EndTime.Sub(s.StartTime).Milliseconds())
	s.Channel <- *s
}

func (s *Span) AddLog(logType string, log string) {
	s.Logs = append(s.Logs, &agent.KeyValue{
		Type:  logType,
		Value: log,
	})
}

func (sid *SID) ToString() string {
	return hex.EncodeToString(sid[:])
}
