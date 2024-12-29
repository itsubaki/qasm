package visitor

import (
	"github.com/itsubaki/q"
	"github.com/oklog/ulid/v2"
)

type Environ struct {
	ID           string
	Version      string
	Qubit        map[string][]q.Qubit
	ClassicalBit map[string][]int64
	Outer        *Environ
}

func NewEnviron() *Environ {
	return &Environ{
		ID:           ulid.Make().String(),
		Qubit:        make(map[string][]q.Qubit),
		ClassicalBit: make(map[string][]int64),
	}
}

func (e *Environ) NewEnclosed() *Environ {
	env := NewEnviron()
	env.Outer = e
	return env
}

func (e *Environ) GetQubit(name string) ([]q.Qubit, bool) {
	if q, ok := e.Qubit[name]; ok {
		return q, true
	}

	if e.Outer != nil {
		return e.Outer.GetQubit(name)
	}

	return nil, false
}

func (e *Environ) GetClassicalBit(name string) ([]int64, bool) {
	if q, ok := e.ClassicalBit[name]; ok {
		return q, true
	}

	if e.Outer != nil {
		return e.Outer.GetClassicalBit(name)
	}

	return nil, false
}
