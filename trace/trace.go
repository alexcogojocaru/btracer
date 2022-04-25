package trace

type ReadTrace interface {
	GenerateID() []byte
}

type Trace struct {
	Encoder Encoder
	TraceID [16]byte
}

func (t *Trace) GenerateID() []byte {
	token := t.Encoder.Compute()
	copy(t.TraceID[:], token)

	return token
}
