package evaluator_test

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleEvaluator() {
	p := &ast.OpenQASM{
		Version: "3.0",
		Includes: []string{
			"\"stdgates.qasm\"",
		},
		Statements: []ast.Stmt{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
					Index: &ast.IndexExpr{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
			},
			&ast.DeclStmt{
				Kind: lexer.BIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
					Index: &ast.IndexExpr{
						Kind:  lexer.INT,
						Value: "2",
					},
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
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.CX,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "q",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "1",
						},
					},
				},
			},
			&ast.ArrowStmt{
				Kind: lexer.ARROW,
				Left: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Target: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
				},
				Right: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
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
					Target: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
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

	if err := e.R.Bit.Println(); err != nil {
		fmt.Println(err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// OPENQASM 3.0;
	// include "stdgates.qasm";
	// qubit q[2];
	// bit c[2];
	// reset q;
	// x q;
	// cx q[0], q[1];
	// measure q -> c;
	// c = measure q;
	//
	// c: 10
	// [10][  2]( 1.0000 0.0000i): 1.0000
}

func ExampleEvaluator_print() {
	p := &ast.OpenQASM{
		Version: "3.0",
		Includes: []string{
			"\"stdgates.qasm\"",
		},
		Statements: []ast.Stmt{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
					Index: &ast.IndexExpr{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.H,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.PrintStmt{
				Kind: lexer.PRINT,
			},
		},
	}

	if err := evaluator.Default().Eval(p); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}
