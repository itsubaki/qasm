package evaluator_test

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleEvaluator() {
	p := &ast.OpenQASM{
		Version: &ast.BasicLit{
			Kind:  lexer.FLOAT,
			Value: "3.0",
		},
		Incls: []ast.Stmt{
			&ast.InclStmt{
				Path: ast.BasicLit{
					Kind:  lexer.STRING,
					Value: "\"stdgates.qasm\"",
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IndexExpr{
						Name: ast.IdentExpr{
							Value: lexer.Tokens[lexer.QUBIT],
						},
						Value: "2",
					},
					Name: ast.IdentExpr{
						Value: "q",
					},
				},
			},
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.BIT,
					Type: &ast.IndexExpr{
						Name: ast.IdentExpr{
							Value: lexer.Tokens[lexer.BIT],
						},
						Value: "2",
					},
					Name: ast.IdentExpr{
						Value: "c",
					},
				},
			},
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Name: lexer.Tokens[lexer.X],
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IndexExpr{
							Name: ast.IdentExpr{
								Value: "q",
							},
							Value: "0",
						},
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.CX,
				Name: lexer.Tokens[lexer.CX],
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IndexExpr{
							Name: ast.IdentExpr{
								Value: "q",
							},
							Value: "0",
						},
						&ast.IndexExpr{
							Name: ast.IdentExpr{
								Value: "q",
							},
							Value: "1",
						},
					},
				},
			},
			&ast.AssignStmt{
				Left: &ast.IdentExpr{
					Value: "c",
				},
				Right: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q",
							},
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

	if err := e.Println(); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// OPENQASM 3.0;
	// include "stdgates.qasm";
	// qubit[2] q;
	// bit[2] c;
	// reset q;
	// x q[0];
	// cx q[0], q[1];
	// c = measure q;
	//
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: 11
}

func ExampleEvaluator_println() {
	p := &ast.OpenQASM{
		Version: &ast.BasicLit{
			Kind:  lexer.FLOAT,
			Value: "3.0",
		},
		Incls: []ast.Stmt{
			&ast.InclStmt{
				Path: ast.BasicLit{
					Kind:  lexer.STRING,
					Value: "\"stdgates.qasm\"",
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IndexExpr{
						Name: ast.IdentExpr{
							Value: lexer.Tokens[lexer.QUBIT],
						},
						Value: "2",
					},
					Name: ast.IdentExpr{
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.H,
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			&ast.PrintStmt{},
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

func ExampleEvaluator_call() {
	p := &ast.OpenQASM{
		Version: &ast.BasicLit{
			Kind:  lexer.FLOAT,
			Value: "3.0",
		},
		Incls: []ast.Stmt{
			&ast.InclStmt{
				Path: ast.BasicLit{
					Kind:  lexer.STRING,
					Value: "\"stdgates.qasm\"",
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
					Name: "bell",
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q0",
							},
							&ast.IdentExpr{
								Value: "q1",
							},
						},
					},
					Body: ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ApplyStmt{
								Kind: lexer.H,
								QArgs: ast.ExprList{
									List: []ast.Expr{
										&ast.IdentExpr{
											Value: "q0",
										},
									},
								},
							},
							&ast.ApplyStmt{
								Kind: lexer.CX,
								QArgs: ast.ExprList{
									List: []ast.Expr{
										&ast.IdentExpr{
											Value: "q0",
										},
										&ast.IdentExpr{
											Value: "q1",
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IndexExpr{
						Name: ast.IdentExpr{
							Value: lexer.Tokens[lexer.QUBIT],
						},
						Value: "2",
					},
					Name: ast.IdentExpr{
						Value: "q",
					},
				},
			},
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Name: "bell",
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: "q",
								},
								Value: "0",
							},
							&ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: "q",
								},
								Value: "1",
							},
						},
					},
				},
			},
			&ast.PrintStmt{},
		},
	}

	if err := evaluator.Default().Eval(p); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}
