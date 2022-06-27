package env

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/ast"
)

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
