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
					Type: ast.IndexExpr{
						Name: ast.IdentExpr{
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
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			"reset q;",
		},
		{
			&ast.ResetStmt{
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
			"reset q[0], q[1];",
		},
		{
			&ast.PrintStmt{},
			"print;",
		},
		{
			&ast.PrintStmt{
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
			"measure q[0], q[1];",
		},
		{
			&ast.ApplyStmt{
				Kind: lexer.X,
				Name: lexer.Tokens[lexer.X],
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
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
					Name: ast.IdentExpr{
						Value: "N",
					},
					Value: &ast.BasicLit{
						Kind:  lexer.INT,
						Value: "15",
					},
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
					Name: ast.IdentExpr{
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
					Name: ast.IdentExpr{
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
			"qubit[2] q;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
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
							&ast.ApplyStmt{
								Kind: lexer.H,
								Name: lexer.Tokens[lexer.H],
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
								Name: lexer.Tokens[lexer.CX],
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
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
					Name: "CX",
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
							&ast.ApplyStmt{
								Kind: lexer.X,
								Name: lexer.Tokens[lexer.X],
								Modifier: []ast.Modifiler{
									{
										Kind: lexer.CTRL,
									},
								},
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
			"gate CX q0, q1 { ctrl @ x q0, q1; }",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
					Name: "CX",
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
							&ast.ApplyStmt{
								Kind: lexer.X,
								Name: lexer.Tokens[lexer.X],
								Modifier: []ast.Modifiler{
									{
										Kind: lexer.CTRL,
										Index: ast.ParenExpr{
											List: ast.ExprList{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  lexer.INT,
														Value: "0",
													},
												},
											},
										},
									},
								},
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
			"gate CX q0, q1 { ctrl(0) @ x q0, q1; }",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
					Name: "CIQFT",
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
							&ast.ApplyStmt{
								Kind: lexer.QFT,
								Name: lexer.Tokens[lexer.QFT],
								Modifier: []ast.Modifiler{
									{
										Kind: lexer.CTRL,
									},
									{
										Kind: lexer.INV,
									},
								},
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
			"gate CIQFT q0, q1 { ctrl @ inv @ qft q0, q1; }",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
