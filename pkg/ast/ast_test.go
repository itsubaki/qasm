package ast_test

import (
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func TestRegStmt(t *testing.T) {
	var cases = []struct {
		in   ast.Stmt
		want string
	}{
		{
			&ast.RegStmt{
				Token: lexer.BIT,
				Name: &ast.Ident{
					Token: lexer.STRING,
					Value: "c",
				},
			},
			"bit c;",
		},
		{
			&ast.RegStmt{
				Token: lexer.BIT,
				Name: &ast.Ident{
					Token: lexer.STRING,
					Value: "c",
				},
				Type: &ast.Array{
					Lbrack: lexer.LBRACKET,
					Rbrack: lexer.RBRACKET,
					Token:  lexer.INT,
					Index: &ast.Ident{
						Token: lexer.INT,
						Value: "2",
					},
				},
			},
			"bit c[2];",
		},
		{
			&ast.RegStmt{
				Token: lexer.QUBIT,
				Name: &ast.Ident{
					Token: lexer.STRING,
					Value: "q",
				},
			},
			"qubit q;",
		},
		{
			&ast.RegStmt{
				Token: lexer.QUBIT,
				Name: &ast.Ident{
					Token: lexer.STRING,
					Value: "q",
				},
				Type: &ast.Array{
					Lbrack: lexer.LBRACKET,
					Rbrack: lexer.RBRACKET,
					Token:  lexer.INT,
					Index: &ast.Ident{
						Token: lexer.INT,
						Value: "2",
					},
				},
			},
			"qubit q[2];",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestResetStmt(t *testing.T) {
	var cases = []struct {
		in   ast.Stmt
		want string
	}{
		{
			&ast.ResetStmt{
				Token: lexer.RESET,
				Target: &ast.RegStmt{
					Token: lexer.QUBIT,
					Name: &ast.Ident{
						Token: lexer.STRING,
						Value: "q",
					},
				},
			},
			"reset q;",
		},
		{
			&ast.ResetStmt{
				Token: lexer.RESET,
				Target: &ast.RegStmt{
					Token: lexer.QUBIT,
					Name: &ast.Ident{
						Token: lexer.STRING,
						Value: "q",
					},
					Type: &ast.Array{
						Lbrack: lexer.LBRACKET,
						Rbrack: lexer.RBRACKET,
						Token:  lexer.INT,
						Index: &ast.Ident{
							Token: lexer.INT,
							Value: "2",
						},
					},
				},
			},
			"reset q[2];",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
