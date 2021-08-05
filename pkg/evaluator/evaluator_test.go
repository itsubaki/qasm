package evaluator_test

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleEvaluator_Eval() {
	p := &ast.OpenQASM{
		Version: 3.0,
		Statements: []ast.Stmt{
			&ast.LetStmt{
				Kind: lexer.QUBIT,
				Index: &ast.IndexExpr{
					LBRACKET: lexer.LBRACKET,
					RBRACKET: lexer.RBRACKET,
					Kind:     lexer.INT,
					Value:    "2",
				},
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.LetStmt{
				Kind: lexer.BIT,
				Index: &ast.IndexExpr{
					LBRACKET: lexer.LBRACKET,
					RBRACKET: lexer.RBRACKET,
					Kind:     lexer.INT,
					Value:    "2",
				},
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Target: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Target: &ast.IdentExpr{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
		},
	}

	fmt.Println(p)

	e := evaluator.Default()
	if err := e.Eval(p); err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range e.Bit {
		for i, vv := range v {
			fmt.Printf("%v[%v]: %v\n", k, i, vv)
		}
	}
	fmt.Println()

	for _, s := range e.Q.State() {
		fmt.Println(s)
	}

	// Output:
	// qubit[2] q;
	// bit[2] c;
	// reset q;
	// x q;
	// c = measure q;
	//
	// c[0]: 1
	// c[1]: 1
	//
	// [11][  3]( 1.0000 0.0000i): 1.0000
}
