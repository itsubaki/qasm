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

func NewQubit() *Qubit {
	return &Qubit{
		Name:  make([]string, 0),
		Value: make(map[string][]q.Qubit),
	}
}

func (qb *Qubit) Add(name string, value []q.Qubit) {
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
		if index > len(out)-1 || index < 0 {
			msg := fmt.Sprintf("index out of range[%v] with length %v", index, len(out))
			panic(msg)
		}

		return append(make([]q.Qubit, 0), out[index]), ok
	default:
		panic(fmt.Sprintf("invalid expr=%#v", a))
	}
}
