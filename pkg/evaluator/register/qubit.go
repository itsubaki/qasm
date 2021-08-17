package register

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
)

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

	index := Index(expr[0].Int(), len(out))
	if index > len(out)-1 || index < 0 {
		return out, fmt.Errorf("index out of range[%v] with length %v", index, len(out))
	}

	return append(make([]q.Qubit, 0), out[index]), nil
}
