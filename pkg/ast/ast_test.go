package ast_test

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleOpenQASM_String() {
	p := &ast.OpenQASM{
		Version: &ast.BasicExpr{
			Kind:  lexer.FLOAT,
			Value: "3.0",
		},
		Incls: []ast.Stmt{
			&ast.InclStmt{
				Path: ast.BasicExpr{
					Kind:  lexer.STRING,
					Value: "\"stdgates.qasm\"",
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IdentExpr{
						Value: lexer.Tokens[lexer.QUBIT],
					},
					Name: ast.IdentExpr{
						Value: "q",
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
		},
	}

	fmt.Println(p)

	// Output:
	// OPENQASM 3.0;
	// include "stdgates.qasm";
	// qubit q;
	// reset q;
}
