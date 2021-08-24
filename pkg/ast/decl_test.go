package ast_test

import (
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func TestDecl(t *testing.T) {
	var cases = []struct {
		in   ast.Decl
		want string
	}{
		{
			&ast.GenConst{
				Name: ast.IdentExpr{
					Value: "N",
				},
				Value: ast.BasicExpr{
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
					Value: lexer.Tokens[lexer.BIT],
				},
				Name: ast.IdentExpr{
					Value: "c",
				},
			},
			"bit c",
		},
		{
			&ast.GenDecl{
				Kind: lexer.QUBIT,
				Type: &ast.IdentExpr{
					Value: lexer.Tokens[lexer.QUBIT],
				},
				Name: ast.IdentExpr{
					Value: "q",
				},
			},
			"qubit q",
		},
		{
			&ast.GenDecl{
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
			"qubit[2] q",
		},
		{
			&ast.GenDecl{
				Kind: lexer.INT,
				Type: &ast.IndexExpr{
					Name: ast.IdentExpr{
						Value: "int",
					},
					Value: "32",
				},
				Name: ast.IdentExpr{
					Value: "a",
				},
			},
			"int[32] a",
		},
		{
			&ast.GateDecl{
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
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.ApplyExpr{
								Kind: lexer.H,
								QArgs: ast.ExprList{
									List: []ast.Expr{
										&ast.IdentExpr{
											Value: "q0",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.ApplyExpr{
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
			"gate bell q0, q1 { h q0; cx q0, q1; }",
		},
		{
			&ast.ParenDecl{
				List: ast.DeclList{
					List: []ast.Decl{
						&ast.GenDecl{
							Kind: lexer.INT,
							Type: &ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: "int",
								},
								Value: "32",
							},
							Name: ast.IdentExpr{
								Value: "a",
							},
						},
						&ast.GenDecl{
							Kind: lexer.INT,
							Type: &ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: "int",
								},
								Value: "32",
							},
							Name: ast.IdentExpr{
								Value: "N",
							},
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
									Name: ast.IdentExpr{
										Value: "int",
									},
									Value: "32",
								},
								Name: ast.IdentExpr{
									Value: "a",
								},
							},
							&ast.GenDecl{
								Kind: lexer.INT,
								Type: &ast.IndexExpr{
									Name: ast.IdentExpr{
										Value: "int",
									},
									Value: "32",
								},
								Name: ast.IdentExpr{
									Value: "N",
								},
							},
						},
					},
				},
				QArgs: ast.DeclList{
					List: []ast.Decl{
						&ast.GenDecl{
							Kind: lexer.QUBIT,
							Type: &ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: lexer.Tokens[lexer.QUBIT],
								},
								Value: "n",
							},
							Name: ast.IdentExpr{
								Value: "r0",
							},
						},
						&ast.GenDecl{
							Kind: lexer.QUBIT,
							Type: &ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: lexer.Tokens[lexer.QUBIT],
								},
								Value: "m",
							},
							Name: ast.IdentExpr{
								Value: "r1",
							},
						},
					},
				},
				Body: ast.BlockStmt{},
				Result: &ast.IndexExpr{
					Name: ast.IdentExpr{
						Value: "bit",
					},
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
