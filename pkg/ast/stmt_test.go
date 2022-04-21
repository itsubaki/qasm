package ast_test

import (
	"fmt"
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
								Name: "q",
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
					Name: "c",
				},
				Right: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Name: "q",
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
								Name: "q",
							},
						},
					},
				},
				Right: &ast.IdentExpr{
					Name: "c",
				},
			},
			"measure q -> c;",
		},
		{
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Name: "q",
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
							Name:  "q",
							Value: "0",
						},
						&ast.IndexExpr{
							Name:  "q",
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
							Name:  "q",
							Value: "0",
						},
						&ast.IndexExpr{
							Name:  "q",
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
								Name: "q",
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
								Name:  "q",
								Value: "0",
							},
							&ast.IndexExpr{
								Name:  "q",
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
							Name: "q",
						},
					},
				},
			},
			"X q;",
		},
		{
			&ast.ExprStmt{
				X: &ast.CallExpr{
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
									Name: "a",
								},
								&ast.IdentExpr{
									Name: "N",
								},
							},
						},
					},
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{
								Name: "r0",
							},
							&ast.IdentExpr{
								Name: "r1",
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
					Name: "N",
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
						Name: lexer.Tokens[lexer.BIT],
					},
					Name: "c",
				},
			},
			"bit c;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IdentExpr{
						Name: lexer.Tokens[lexer.QUBIT],
					},
					Name: "q",
				},
			},
			"qubit q;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IndexExpr{
						Name:  lexer.Tokens[lexer.QUBIT],
						Value: "2",
					},
					Name: "q",
				},
			},
			"qubit[2] q;",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.INT,
					Type: &ast.IndexExpr{
						Name:  "int",
						Value: "32",
					},
					Name: "a",
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
			},
			"gate bell q0, q1 { H q0; cx q0, q1; }",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
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
					Body: ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ApplyStmt{
								Kind: lexer.X,
								Name: lexer.Tokens[lexer.X],
								Modifier: []ast.Modifier{
									{
										Kind: lexer.CTRL,
									},
								},
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
			"gate cx q0, q1 { ctrl @ X q0, q1; }",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
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
					Body: ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ApplyStmt{
								Kind: lexer.X,
								Name: lexer.Tokens[lexer.X],
								Modifier: []ast.Modifier{
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
			"gate cx q0, q1 { ctrl(0) @ X q0, q1; }",
		},
		{
			&ast.DeclStmt{
				Decl: &ast.GateDecl{
					Name: "ciqft",
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
								Kind: lexer.QFT,
								Name: lexer.Tokens[lexer.QFT],
								Modifier: []ast.Modifier{
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
			"gate ciqft q0, q1 { ctrl @ inv @ QFT q0, q1; }",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func ExampleBlockStmt_Reverse() {
	block := &ast.BlockStmt{
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
	}

	r := block.Reverse()

	fmt.Println(block)
	fmt.Println(&r)

	// Output:
	// { H q0; cx q0, q1; }
	// { cx q0, q1; H q0; }
}
