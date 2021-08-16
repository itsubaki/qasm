package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
)

type Gate map[string]ast.GateStmt

type Const map[string]int

type Qubit struct {
	Name  []string
	Value map[string][]q.Qubit
}

func (qb *Qubit) Add(name string, value []q.Qubit) {
	qb.Name = append(qb.Name, name)
	qb.Value[name] = value
}

func (qb *Qubit) Exists(name string) bool {
	_, ok := qb.Value[name]
	return ok
}

func (qb *Qubit) Get(name string, expr ...*ast.IndexExpr) ([]q.Qubit, error) {
	out, ok := qb.Value[name]
	if !ok {
		return nil, fmt.Errorf("IDENT=%v not found", name)
	}
	if len(expr) == 0 {
		return out, nil
	}
	if expr[0] == nil {
		return out, nil
	}

	index := expr[0].Int()
	if index < 0 {
		index = len(out) + index
	}

	if index > len(out)-1 || index < 0 {
		return out, fmt.Errorf("index out of range[%v] with length %v", index, len(out))
	}

	return append(make([]q.Qubit, 0), out[index]), nil
}

type Bit struct {
	Name  []string
	Value map[string][]int
}

func (b *Bit) Add(name string, value []int) {
	b.Name = append(b.Name, name)
	b.Value[name] = value
}

func (b *Bit) Exists(name string) bool {
	_, ok := b.Value[name]
	return ok
}

func (b *Bit) Get(name string, expr ...*ast.IndexExpr) ([]int, error) {
	out, ok := b.Value[name]
	if !ok {
		return nil, fmt.Errorf("IDENT=%v not found", name)
	}
	if len(expr) == 0 {
		return out, nil
	}
	if expr[0] == nil {
		return out, nil
	}

	index := expr[0].Int()
	if index < 0 {
		index = len(out) + index
	}

	if index > len(out)-1 || index < 0 {
		return out, fmt.Errorf("index out of range[%v] with length %v", index, len(out))
	}

	return append(make([]int, 0), out[index]), nil
}

func (b *Bit) Println() error {
	for _, n := range b.Name {
		fmt.Printf("%v: ", n)

		c, err := b.Get(n)
		if err != nil {
			return fmt.Errorf("get bit=%v: %v", n, err)
		}

		for _, v := range c {
			fmt.Printf("%v", v)
		}

		fmt.Println()
	}

	return nil
}
