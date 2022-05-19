package env

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
)

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

func (b *Bit) String() string {
	return fmt.Sprintf("%v, %v", b.Name, b.Value)
}
