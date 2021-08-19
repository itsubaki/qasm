package ast_test

import (
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func TestDeclStmt(t *testing.T) {
	var cases = []struct {
		in   ast.Stmt
		want string
	}{
		{
			&ast.DeclStmt{
				Decl: &ast.GenConst{
					Name: &ast.IdentExpr{
						Value: "N",
					},
					Value: "15",
				},
			},
			"const N = 15;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.BIT,
					Type: &ast.IdentExpr{
						Value: "bit",
					},
					Name: &ast.IdentExpr{
						Value: "c",
					},
				},
			},
			"bit c;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IdentExpr{
						Value: "qubit",
					},
					Name: &ast.IdentExpr{
						Value: "q",
					},
				},
			},
			"qubit q;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IndexExpr{
						Name: &ast.IdentExpr{
							Value: "qubit",
						},
						Value: "2",
					},
					Name: &ast.IdentExpr{
						Value: "q",
					},
				},
			},
			"qubit[2] q;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.INT,
					Type: &ast.IndexExpr{
						Name: &ast.IdentExpr{
							Value: "int",
						},
						Value: "32",
					},
					Name: &ast.IdentExpr{
						Value: "a",
					},
				},
			},
			"int[32] a;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.FuncDecl{
					Kind: lexer.GATE,
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
			},
			"gate bell q0, q1 { h q0; cx q0, q1; };",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
