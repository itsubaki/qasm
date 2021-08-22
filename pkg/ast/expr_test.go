package ast_test

import (
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func TestExpr(t *testing.T) {
	var cases = []struct {
		in   ast.Expr
		want string
	}{
		{
			&ast.ArrayExpr{
				Type: &ast.IndexExpr{
					Name: &ast.IdentExpr{
						Value: "int",
					},
					Value: "32",
				},
				Name: "a",
			},
			"int[32] a",
		},
		{
			&ast.ResetExpr{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			"reset q",
		},
		{
			&ast.ResetExpr{
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
			"reset q[0], q[1]",
		},
		{
			&ast.PrintExpr{},
			"print",
		},
		{
			&ast.PrintExpr{
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
			"print q[0], q[1]",
		},
		{
			&ast.MeasureExpr{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			"measure q",
		},
		{
			&ast.MeasureExpr{
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
			"measure q[0], q[1]",
		},
		{
			&ast.ApplyExpr{
				Kind: lexer.X,
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			"x q",
		},
		{
			&ast.CallExpr{
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
			"bell q0, q1",
		},
		{
			&ast.CallExpr{
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
			"shor(a, N) r0, r1",
		},
		{
			&ast.ParenExpr{
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
			"(a, N)",
		},
		{
			&ast.ParenExpr{
				List: ast.ExprList{
					List: []ast.Expr{
						&ast.ArrayExpr{
							Type: &ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "int",
								},
								Value: "32",
							},
							Name: "a",
						},
						&ast.ArrayExpr{
							Type: &ast.IndexExpr{
								Name: &ast.IdentExpr{
									Value: "int",
								},
								Value: "32",
							},
							Name: "N",
						},
					},
				},
			},
			"(int[32] a, int[32] N)",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
