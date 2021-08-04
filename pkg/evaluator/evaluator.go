package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
)

type Evaluator struct {
	qsim *q.Q
}

func New() *Evaluator {
	return &Evaluator{
		qsim: q.New(),
	}
}

func (e *Evaluator) Eval(p *ast.Program) error {
	for _, s := range p.Statements {
		switch s.(type) {
		case *ast.LetStmt:
			fmt.Println(s)
		case *ast.ResetStmt:
			fmt.Println(s)
		case *ast.ApplyStmt:
			fmt.Println(s)
		case *ast.AssignStmt:
			fmt.Println(s)
		default:
			return fmt.Errorf("invalid stmt=%s", s)
		}
	}

	return nil
}
