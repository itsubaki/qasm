package evaluator_test

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleEvaluator_Eval() {
	p := &ast.Program{
		Statements: []ast.Stmt{
			&ast.LetStmt{
				Kind: lexer.QUBIT,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
					Index: &ast.Index{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				Name: []ast.Ident{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.MeasureStmt{
				Kind: lexer.MEASURE,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
		},
	}

	fmt.Println(p)

	e := evaluator.Default()
	if err := e.Eval(p); err != nil {
		fmt.Println(err)
	}

	for _, s := range e.QSim.State() {
		fmt.Println(s)
	}

	// Output:
	// qubit q[2];
	// reset q;
	// x q;
	// measure q;
	//
	// [11][  3]( 1.0000 0.0000i): 1.0000
}
