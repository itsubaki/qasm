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
				Type: ast.IndexExpr{
					Name: ast.IdentExpr{
						Value: "int",
					},
					Value: "32",
				},
				Name: "a",
			},
			"int[32] a",
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
			"measure q[0], q[1]",
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
				Name: "bell",
				Modifier: []ast.Modifier{
					{
						Kind: lexer.POW,
						Index: ast.ParenExpr{
							List: ast.ExprList{
								List: []ast.Expr{
									&ast.UnaryExpr{
										Kind: lexer.MINUS,
										Value: &ast.BasicLit{
											Kind:  lexer.INT,
											Value: "2",
										},
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
			"pow(-2) @ bell q0, q1",
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
							Type: ast.IndexExpr{
								Name: ast.IdentExpr{
									Value: "int",
								},
								Value: "32",
							},
							Name: "a",
						},
						&ast.ArrayExpr{
							Type: ast.IndexExpr{
								Name: ast.IdentExpr{
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
		{
			&ast.InfixExpr{
				Kind: lexer.PLUS,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "1",
				},
			},
			"3 + 1",
		},
		{

			&ast.ParenExpr{
				List: ast.ExprList{
					List: []ast.Expr{
						&ast.InfixExpr{
							Kind: lexer.PLUS,
							Left: &ast.BasicLit{
								Kind:  lexer.INT,
								Value: "3",
							},
							Right: &ast.ParenExpr{
								List: ast.ExprList{
									List: []ast.Expr{
										&ast.InfixExpr{
											Kind: lexer.MINUS,
											Left: &ast.BasicLit{
												Kind:  lexer.INT,
												Value: "5",
											},
											Right: &ast.IdentExpr{
												Value: "a",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"(3 + (5 - a))",
		},
		{
			&ast.ParenExpr{
				List: ast.ExprList{
					List: []ast.Expr{
						&ast.InfixExpr{
							Kind: lexer.PLUS,
							Left: &ast.BasicLit{
								Kind:  lexer.INT,
								Value: "3",
							},
							Right: &ast.BasicLit{
								Kind:  lexer.INT,
								Value: "1",
							},
						},
					},
				},
			},
			"(3 + 1)",
		},
		{
			&ast.UnaryExpr{
				Kind: lexer.MINUS,
				Value: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "1",
				},
			},
			"-1",
		},
		{
			&ast.UnaryExpr{
				Kind: lexer.MINUS,
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "1.0",
				},
			},
			"-1.0",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestModPow(t *testing.T) {
	mod := []ast.Modifier{
		{Kind: lexer.POW},
		{Kind: lexer.INV},
		{Kind: lexer.CTRL},
		{Kind: lexer.NEGCTRL},
	}

	pow := ast.ModPow(mod)
	for _, p := range pow {
		if p.Kind != lexer.POW {
			t.Errorf("invalid kind=%v", p.Kind)
		}
	}
}

func TestModInv(t *testing.T) {
	mod := []ast.Modifier{
		{Kind: lexer.POW},
		{Kind: lexer.INV},
		{Kind: lexer.CTRL},
		{Kind: lexer.NEGCTRL},
	}

	inv := ast.ModInv(mod)
	for _, p := range inv {
		if p.Kind != lexer.INV {
			t.Errorf("invalid kind=%v", p.Kind)
		}
	}
}

func TestModCtrl(t *testing.T) {
	mod := []ast.Modifier{
		{Kind: lexer.POW},
		{Kind: lexer.INV},
		{Kind: lexer.CTRL},
		{Kind: lexer.NEGCTRL},
	}

	ctrl := ast.ModCtrl(mod)
	for _, p := range ctrl {
		if p.Kind != lexer.CTRL && p.Kind != lexer.NEGCTRL {
			t.Errorf("invalid kind=%v", p.Kind)
		}
	}
}
