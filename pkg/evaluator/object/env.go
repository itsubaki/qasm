package object

import (
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
	name := ast.Ident(n)
	b.Name = append(b.Name, name)
	b.Value[name] = value
}

func (b *Bit) Get(a ast.Expr) ([]int64, bool) {
	switch x := a.(type) {
	case *ast.IdentExpr:
		out, ok := b.Value[x.Value]
		return out, ok
	case *ast.ArrayExpr:
		out, ok := b.Value[x.Name]
		return out, ok
	case *ast.IndexExpr:
		out, ok := b.Value[x.Name.Value]
		index := Index(x.Int(), len(out))
		return append(make([]int64, 0), out[index]), ok
	}

	panic(fmt.Sprintf("invalid expr=%#v", a))
}

type Qubit struct {
	Name  []string
	Value map[string][]q.Qubit
}

func (qb *Qubit) Add(n ast.Node, value []q.Qubit) {
	name := ast.Ident(n)
	qb.Name = append(qb.Name, name)
	qb.Value[name] = value
}

func (qb *Qubit) Get(a ast.Expr) ([]q.Qubit, bool) {
	switch x := a.(type) {
	case *ast.IdentExpr:
		out, ok := qb.Value[x.Value]
		return out, ok
	case *ast.ArrayExpr:
		out, ok := qb.Value[x.Name]
		return out, ok
	case *ast.IndexExpr:
		out, ok := qb.Value[x.Name.Value]
		index := Index(x.Int(), len(out))
		return append(make([]q.Qubit, 0), out[index]), ok
	}

	panic(fmt.Sprintf("invalid expr=%#v", a))
}

func Index(index, len int) int {
	out := index
	if index < 0 {
		out = len + index
	}

	if out > len || out < 0 {
		msg := fmt.Sprintf("index out of range[%v] with length %v", out, len)
		panic(msg)
	}

	return out
}
