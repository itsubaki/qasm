package object

import (
	"bytes"
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
)

type Environment struct {
	Bit   *Bit
	Qubit *Qubit
	Const Const
	Func  Func
	Outer *Environment
}

func (e *Environment) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("const: %v, ", e.Const))
	buf.WriteString(fmt.Sprintf("bit: %v, ", e.Bit.Value))
	buf.WriteString(fmt.Sprintf("qubit: %v, ", e.Qubit.Value))
	buf.WriteString(fmt.Sprintf("func: %v, ", e.Func))

	return buf.String()
}

type Func map[string]ast.Decl

func NewEnvironment() *Environment {
	return &Environment{
		Bit: &Bit{
			Name:  make([]string, 0),
			Value: make(map[string][]int64),
		},
		Qubit: &Qubit{
			Name:  make([]string, 0),
			Value: make(map[string][]q.Qubit),
		},
		Const: make(map[string]Object),
		Func:  make(map[string]ast.Decl),
		Outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer
	env.Func = outer.Func
	env.Const = outer.Const

	return env
}

type Const map[string]Object

type Bit struct {
	Name  []string
	Value map[string][]int64
}

func (b *Bit) Add(n ast.Node, value []int64) {
	name := ast.Must(ast.Ident(n))
	b.Name = append(b.Name, name)
	b.Value[name] = value
}

func (b *Bit) Get(a ast.Expr) ([]int64, bool) {
	switch x := a.(type) {
	case *ast.IdentExpr:
		out, ok := b.Value[x.Name]
		return out, ok
	case *ast.IndexExpr:
		out, ok := b.Value[x.Name]
		idx, err := index(x.Int(), len(out))
		if err != nil {
			return nil, false
		}

		return []int64{out[idx]}, ok
	}

	return nil, false
}

type Qubit struct {
	Name  []string
	Value map[string][]q.Qubit
}

func (qb *Qubit) Add(n ast.Node, value []q.Qubit) {
	name := ast.Must(ast.Ident(n))
	qb.Name = append(qb.Name, name)
	qb.Value[name] = value
}

func (qb *Qubit) Get(a ast.Expr) ([]q.Qubit, bool) {
	switch x := a.(type) {
	case *ast.IdentExpr:
		out, ok := qb.Value[x.Name]
		return out, ok
	case *ast.IndexExpr:
		out, ok := qb.Value[x.Name]
		idx, err := index(x.Int(), len(out))
		if err != nil {
			return nil, false
		}

		return []q.Qubit{out[idx]}, ok
	}

	return nil, false
}

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
