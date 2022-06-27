package ast_test

import (
	"testing"

	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/lexer"
)

func TestDecl(t *testing.T) {
	var cases = []struct {
		in   ast.Decl
		want string
	}{
		{
			&ast.VersionDecl{
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3.0",
				},
			},
			"OPENQASM 3.0",
		},
		{
			&ast.GenConst{
				Name: "N",
				Value: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "15",
				},
			},
			"const N = 15",
		},
		{
			&ast.GenDecl{
				Kind: lexer.BIT,
				Type: &ast.IdentExpr{
					Name: lexer.Tokens[lexer.BIT],
				},
				Name: "c",
			},
			"bit c",
		},
		{
			&ast.GenDecl{
				Kind: lexer.QUBIT,
				Type: &ast.IdentExpr{
					Name: lexer.Tokens[lexer.QUBIT],
				},
				Name: "q",
			},
			"qubit q",
		},
		{
			&ast.GenDecl{
				Kind: lexer.QUBIT,
				Type: &ast.IndexExpr{
					Name:  lexer.Tokens[lexer.QUBIT],
					Value: "2",
				},
				Name: "q",
			},
			"qubit[2] q",
		},
		{
			&ast.GenDecl{
				Kind: lexer.INT,
				Type: &ast.IndexExpr{
					Name:  "int",
					Value: "32",
				},
				Name: "a",
			},
			"int[32] a",
		},
		{
			&ast.GateDecl{
				Name: "bell",
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Name: "q0",
						},
						&ast.IdentExpr{
							Name: "q1",
						},
					},
				},
				Body: ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ApplyStmt{
							Kind: lexer.H,
							Name: lexer.Tokens[lexer.H],
							QArgs: ast.ExprList{
								List: []ast.Expr{
									&ast.IdentExpr{
										Name: "q0",
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Name: "cx",
								QArgs: ast.ExprList{
									List: []ast.Expr{
										&ast.IdentExpr{
											Name: "q0",
										},
										&ast.IdentExpr{
											Name: "q1",
										},
									},
								},
							},
						},
					},
				},
			},
			"gate bell q0, q1 { H q0; cx q0, q1; }",
		},
		{
			&ast.ParenDecl{
				List: ast.DeclList{
					List: []ast.Decl{
						&ast.GenDecl{
							Kind: lexer.INT,
							Type: &ast.IndexExpr{
								Name:  "int",
								Value: "32",
							},
							Name: "a",
						},
						&ast.GenDecl{
							Kind: lexer.INT,
							Type: &ast.IndexExpr{
								Name:  "int",
								Value: "32",
							},
							Name: "N",
						},
					},
				},
			},
			"(int[32] a, int[32] N)",
		},
		{
			&ast.FuncDecl{
				Name: "shor",
				Params: ast.ParenDecl{
					List: ast.DeclList{
						List: []ast.Decl{
							&ast.GenDecl{
								Kind: lexer.INT,
								Type: &ast.IndexExpr{
									Name:  "int",
									Value: "32",
								},
								Name: "a",
							},
							&ast.GenDecl{
								Kind: lexer.INT,
								Type: &ast.IndexExpr{
									Name:  "int",
									Value: "32",
								},
								Name: "N",
							},
						},
					},
				},
				QArgs: ast.DeclList{
					List: []ast.Decl{
						&ast.GenDecl{
							Kind: lexer.QUBIT,
							Type: &ast.IndexExpr{
								Name:  lexer.Tokens[lexer.QUBIT],
								Value: "n",
							},
							Name: "r0",
						},
						&ast.GenDecl{
							Kind: lexer.QUBIT,
							Type: &ast.IndexExpr{
								Name:  lexer.Tokens[lexer.QUBIT],
								Value: "m",
							},
							Name: "r1",
						},
					},
				},
				Body: ast.BlockStmt{},
				Result: &ast.IndexExpr{
					Name:  "bit",
					Value: "n",
				},
			},
			"def shor(int[32] a, int[32] N) qubit[n] r0, qubit[m] r1 -> bit[n] { }",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestGenDecl(t *testing.T) {
	var cases = []struct {
		in   ast.GenDecl
		want int
	}{
		{
			in: ast.GenDecl{
				Kind: lexer.QUBIT,
				Type: &ast.IndexExpr{
					Name:  lexer.Tokens[lexer.QUBIT],
					Value: "2",
				},
				Name: "q",
			},
			want: 2,
		},
		{
			in: ast.GenDecl{
				Kind: lexer.QUBIT,
				Type: &ast.IdentExpr{
					Name: lexer.Tokens[lexer.QUBIT],
				},
				Name: "q",
			},
			want: 1,
		},
	}

	for _, c := range cases {
		got := c.in.Size()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDeclLiteral(t *testing.T) {
	var cases = []struct {
		in   ast.Decl
		want string
	}{
		{&ast.BadDecl{}, ""},
		{&ast.ParenDecl{}, lexer.Tokens[lexer.LPAREN]},
		{&ast.GenDecl{Kind: lexer.QUBIT}, lexer.Tokens[lexer.QUBIT]},
	}

	for _, c := range cases {
		got := c.in.Literal()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestBadDecl(t *testing.T) {
	d := &ast.BadDecl{}

	if len(d.String()) > 0 {
		t.Errorf("invalid string= %v", d.String())
	}

	if len(d.Literal()) > 0 {
		t.Errorf("invalid literal= %v", d.Literal())
	}
}
