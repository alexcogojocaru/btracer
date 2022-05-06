package trace

import "encoding/hex"

type ReadTrace interface {
	GenerateID() []byte
}

const DEFAULT_TRACE_BYTES_SIZE = 16

var NullTraceID TID

type TID [DEFAULT_TRACE_BYTES_SIZE]byte

func (tid *TID) ToString() string {
	return hex.EncodeToString(tid[:])
}

type Trace struct {
	TraceID TID
}
