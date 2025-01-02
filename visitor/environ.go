package visitor

import (
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/oklog/ulid/v2"
)

type Environ struct {
	ID           string
	Version      string
	Qubit        map[string][]q.Qubit
	ClassicalBit map[string][]int64
	Variable     map[string]interface{}
	Gate         map[string]Gate
	Outer        *Environ
}

func NewEnviron() *Environ {
	return &Environ{
		ID:           ulid.Make().String(),
		Qubit:        make(map[string][]q.Qubit),
		ClassicalBit: make(map[string][]int64),
		Variable:     make(map[string]interface{}),
		Gate:         make(map[string]Gate),
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

func (e *Environ) GetVariable(name string) (interface{}, bool) {
	if v, ok := e.Variable[name]; ok {
		return v, true
	}

	if e.Outer != nil {
		return e.Outer.GetVariable(name)
	}

	return nil, false
}

func (e *Environ) GetGate(name string) (Gate, bool) {
	if g, ok := e.Gate[name]; ok {
		return g, true
	}

	if e.Outer != nil {
		return e.Outer.GetGate(name)
	}

	return Gate{}, false
}

type Gate struct {
	Name   string
	Params []string
	QArgs  []string
	Body   []parser.IGateCallStatementContext
}
