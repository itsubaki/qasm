package ast_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/lexer"
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
						Name: lexer.Tokens[lexer.QUBIT],
					},
					Name: "q",
				},
			},
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Name: "q",
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
		in     interface{}
		want   string
		hasErr bool
	}{
		{&ast.IdentExpr{Name: "ident"}, "ident", false},
		{&ast.IndexExpr{Name: "index"}, "index", false},
		{&ast.GenDecl{Name: "gendecl"}, "gendecl", false},
		{&ast.GenConst{Name: "genconst"}, "genconst", false},
		{&ast.GateDecl{Name: "gatedecl"}, "gatedecl", false},
		{&ast.SubroutineDecl{Name: "subroutinedecl"}, "subroutinedecl", false},
		{&ast.BasicLit{Value: "basic"}, "basic", false},
		{"foobar", "", true},
	}

	for _, c := range cases {
		got, err := ast.Ident(c.in)
		if (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}

		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
