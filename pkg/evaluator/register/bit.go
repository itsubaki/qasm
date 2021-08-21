package register

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
)

type Bit struct {
	Name  []string
	Value map[string][]int
}

func (b *Bit) Add(name string, value []int) {
	b.Name = append(b.Name, name)
	b.Value[name] = value
}

func (b *Bit) Get(a ast.Expr) ([]int, bool) {
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
		if index > len(out)-1 || index < 0 {
			msg := fmt.Sprintf("index out of range[%v] with length %v", index, len(out))
			panic(msg)
		}

		return append(make([]int, 0), out[index]), ok
	default:
		panic(fmt.Sprintf("invalid expr=%#v", a))
	}
}

func (b *Bit) Println() error {
	for _, n := range b.Name {
		fmt.Printf("%v: ", n)

		c, ok := b.Get(&ast.IdentExpr{Value: n})
		if !ok {
			return fmt.Errorf("name=%v not found", n)
		}

		for _, v := range c {
			fmt.Printf("%v", v)
		}

		fmt.Println()
	}

	return nil
}
