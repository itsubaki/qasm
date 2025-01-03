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
	Subroutine   map[string]Subroutine
	Outer        *Environ
}

func NewEnviron() *Environ {
	return &Environ{
		ID:           ulid.Make().String(),
		Qubit:        make(map[string][]q.Qubit),
		ClassicalBit: make(map[string][]int64),
		Variable:     make(map[string]interface{}),
		Gate:         make(map[string]Gate),
		Subroutine:   make(map[string]Subroutine),
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

func (e *Environ) GetSubroutine(name string) (Subroutine, bool) {
	if s, ok := e.Subroutine[name]; ok {
		return s, true
	}

	if e.Outer != nil {
		return e.Outer.GetSubroutine(name)
	}

	return Subroutine{}, false
}

type Gate struct {
	Name   string
	Params []string
	QArgs  []string
	Body   []*parser.GateCallStatementContext
}

type Subroutine struct {
	Name string
}

func flatten(qargs [][]q.Qubit) []q.Qubit {
	var flat []q.Qubit
	for _, q := range qargs {
		flat = append(flat, q...)
	}

	return flat
}
