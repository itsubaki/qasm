package ast_test

import (
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func TestStmt(t *testing.T) {
	var cases = []struct {
		in   ast.Stmt
		want string
	}{
		{
			&ast.ReturnStmt{
				Result: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q",
							},
						},
					},
				},
			},
			"return measure q;",
		},
		{
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
			"c = measure q;",
		},
		{
			&ast.ArrowStmt{
				Left: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q",
							},
						},
					},
				},
				Right: &ast.IdentExpr{
					Value: "c",
				},
			},
			"measure q -> c;",
		},
		{
			&ast.ExprStmt{
				X: &ast.ArrayExpr{
					Type: &ast.IndexExpr{
						Name: &ast.IdentExpr{
							Value: "int",
						},
						Value: "32",
					},
					Name: "a",
				},
			},
			"int[32] a;",
		},
		{
			&ast.ExprStmt{
				X: &ast.ResetExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q",
							},
						},
					},
				},
			},
			"reset q;",
		},
		{
			&ast.ExprStmt{
				X: &ast.ResetExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "q",
								},
								Value: "0",
							},
							&ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "q",
								},
								Value: "1",
							},
						},
					},
				},
			},
			"reset q[0], q[1];",
		},
		{
			&ast.ExprStmt{
				X: &ast.PrintExpr{},
			},
			"print;",
		},
		{
			&ast.ExprStmt{
				X: &ast.PrintExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "q",
								},
								Value: "0",
							},
							&ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "q",
								},
								Value: "1",
							},
						},
					},
				},
			},
			"print q[0], q[1];",
		},
		{
			&ast.ExprStmt{
				X: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q",
							},
						},
					},
				},
			},
			"measure q;",
		},
		{
			&ast.ExprStmt{
				X: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "q",
								},
								Value: "0",
							},
							&ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "q",
								},
								Value: "1",
							},
						},
					},
				},
			},
			"measure q[0], q[1];",
		},
		{
			&ast.ExprStmt{
				X: &ast.ApplyExpr{
					Kind: lexer.X,
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "q",
							},
						},
					},
				},
			},
			"x q;",
		},
		{
			&ast.ExprStmt{
				X: &ast.CallExpr{
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
				},
			},
			"bell q0, q1;",
		},
		{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Name: "shor",
					Params: ast.ParenExpr{
						List: ast.ExprList{
							List: []ast.Expr{
								&ast.IdentExpr{
									Value: "a",
								},
								&ast.IdentExpr{
									Value: "N",
								},
							},
						},
					},
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Value: "r0",
							},
							&ast.IdentExpr{
								Value: "r1",
							},
						},
					},
				},
			},
			"shor(a, N) r0, r1;",
		},
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
						Value: lexer.Tokens[lexer.BIT],
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
						Value: lexer.Tokens[lexer.QUBIT],
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
							Value: lexer.Tokens[lexer.QUBIT],
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
			"gate bell q0, q1 { h q0; cx q0, q1; }",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
