package trace

type ReadTrace interface {
	GenerateID() []byte
}

type TID [16]byte

type Trace struct {
	Encoder Encoder
	TraceID TID
}

func (t *Trace) GenerateID() []byte {
	token := t.Encoder.Compute()
	copy(t.TraceID[:], token)

	return token
}
