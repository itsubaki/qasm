package ast_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleOpenQASM_String() {
	p := &ast.OpenQASM{
		Version: &ast.DeclStmt{
			Decl: &ast.VersionDecl{
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3.0",
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.InclStmt{
				Path: ast.BasicLit{
					Kind:  lexer.STRING,
					Value: "\"stdgates.qasm\"",
				},
			},
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

func TestIdent(t *testing.T) {
	var cases = []struct {
		in   interface{}
		want string
	}{
		{&ast.IdentExpr{Value: "ident"}, "ident"},
		{&ast.IndexExpr{Name: ast.IdentExpr{Value: "index"}}, "index"},
		{&ast.ArrayExpr{Name: "array"}, "array"},
		{&ast.GenDecl{Name: ast.IdentExpr{Value: "gendecl"}}, "gendecl"},
		{&ast.GenConst{Name: ast.IdentExpr{Value: "genconst"}}, "genconst"},
		{&ast.GateDecl{Name: "gatedecl"}, "gatedecl"},
		{&ast.FuncDecl{Name: "funcdecl"}, "funcdecl"},
		{&ast.BasicLit{Value: "basic"}, "basic"},
	}

	for _, c := range cases {
		got := ast.Ident(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
