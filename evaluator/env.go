package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/evaluator/object"
)

type Environ struct {
	Outer     *Environ
	Bit       *Bit
	Qubit     *Qubit
	Const     Const
	GateDef   GateDef
	Modifier  []ast.Modifier
	Decl      []ast.Decl
	CtrlQArgs []ast.Expr
}

// NewEnviron returns a new environment.
func NewEnviron() *Environ {
	return &Environ{
		Outer:     nil,
		Bit:       NewBit(),
		Qubit:     NewQubit(),
		GateDef:   make(map[string]ast.Decl),
		Const:     make(map[string]object.Object),
		Modifier:  make([]ast.Modifier, 0),
		Decl:      make([]ast.Decl, 0),
		CtrlQArgs: nil,
	}
}

// NewEnclosed returns a new environment that encloses the outer environment.
func (e *Environ) NewEnclosed(decl ast.Decl, mod []ast.Modifier) *Environ {
	return &Environ{
		Outer:    e,
		Bit:      NewBit(),
		Qubit:    NewQubit(),
		GateDef:  e.GateDef,
		Const:    e.Const,
		Modifier: mod,
		Decl:     append(e.Decl, decl),
	}
}

func (e *Environ) String() string {
	return fmt.Sprintf("gatedef: %v, const: %v, bit: %v, qubit: %v, modifier: %v, decl: %v", e.GateDef, e.Const, e.Bit, e.Qubit, e.Modifier, e.Decl)
}

type GateDef map[string]ast.Decl

type Const map[string]object.Object

type Bit struct {
	Name  []string
	Value map[string][]int64
}

func NewBit() *Bit {
	return &Bit{
		Name:  make([]string, 0),
		Value: make(map[string][]int64),
	}
}

func (b *Bit) Add(n ast.Node, value []int64) {
	name := ast.Must(ast.Ident(n))
	b.Name = append(b.Name, name)
	b.Value[name] = value
}

// Get returns a value of bit.
// if a is ident, return all values.
// if a is index, return a value of index.
func (b *Bit) Get(a ast.Expr) ([]int64, bool) {
	switch x := a.(type) {
	case *ast.IdentExpr:
		out, ok := b.Value[x.Name]
		return out, ok
	case *ast.IndexExpr:
		out, ok := b.Value[x.Name]
		if !ok {
			return nil, false
		}

		idx, err := index(x.Int(), len(out))
		if err != nil {
			return nil, false
		}

		return []int64{out[idx]}, true
	}

	return nil, false
}

func (b *Bit) String() string {
	return fmt.Sprintf("%v, %v", b.Name, b.Value)
}

type Qubit struct {
	Name  []string
	Value map[string][]q.Qubit
}

func NewQubit() *Qubit {
	return &Qubit{
		Name:  make([]string, 0),
		Value: make(map[string][]q.Qubit),
	}
}

func (qb *Qubit) Add(n ast.Node, value []q.Qubit) {
	name := ast.Must(ast.Ident(n))
	qb.Name = append(qb.Name, name)
	qb.Value[name] = value
}

// Get returns a value of qubit.
// if a is ident, return all values.
// if a is index, return a value of index.
func (qb *Qubit) Get(a ast.Expr) ([]q.Qubit, bool) {
	switch x := a.(type) {
	case *ast.IdentExpr:
		out, ok := qb.Value[x.Name]
		return out, ok
	case *ast.IndexExpr:
		out, ok := qb.Value[x.Name]
		if !ok {
			return nil, false
		}

		idx, err := index(x.Int(), len(out))
		if err != nil {
			return nil, false
		}

		return []q.Qubit{out[idx]}, true
	}

	return nil, false
}

// All returns all values of qubit.
func (qb *Qubit) All() []q.Qubit {
	out := make([]q.Qubit, 0)
	for _, n := range qb.Name {
		qb, _ := qb.Get(&ast.IdentExpr{Name: n}) // no error
		out = append(out, qb...)
	}

	return out
}

func (qb *Qubit) String() string {
	return fmt.Sprintf("%v, %v", qb.Name, qb.Value)
}

// q[-1] -> q[len(q)-1]
func index(idx, len int) (int, error) {
	out := idx
	if idx < 0 {
		out = len + idx
	}

	if out > len || out < 0 {
		return 0, fmt.Errorf("index out of range[%v] with length %v", out, len)
	}

	return out, nil
}
