package evaluator_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleEvaluator() {
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
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.BIT,
					Type: &ast.IndexExpr{
						Name: ast.IdentExpr{
							Value: lexer.Tokens[lexer.BIT],
						},
						Value: "2",
					},
					Name: ast.IdentExpr{
						Value: "c",
					},
				},
			},
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Name: lexer.Tokens[lexer.X],
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IndexExpr{
							Name: ast.IdentExpr{
								Value: "q",
							},
							Value: "0",
						},
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.CX,
				Name: lexer.Tokens[lexer.CX],
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
		},
	}

	fmt.Println(p)

	e := evaluator.Default()
	if err := e.Eval(p); err != nil {
		fmt.Println(err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// OPENQASM 3.0;
	// qubit[2] q;
	// bit[2] c;
	// reset q;
	// X q[0];
	// CX q[0], q[1];
	// c = measure q;
	//
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: 11
}

func ExampleEvaluator_call() {
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
					Body: ast.BlockStmt{
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
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Name: "bell",
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
			&ast.PrintStmt{},
		},
	}

	if err := evaluator.Default().Eval(p); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func TestEvalExpr(t *testing.T) {
	var cases = []struct {
		in   ast.Expr
		want object.Object
	}{
		{
			in: &ast.BasicLit{
				Kind:  lexer.INT,
				Value: "3",
			},
			want: &object.Int{
				Value: 3,
			},
		},
		{
			in: &ast.BasicLit{
				Kind:  lexer.PI,
				Value: "pi",
			},
			want: &object.Float{
				Value: 3.141592653589793,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.PLUS,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "11",
				},
			},
			want: &object.Int{
				Value: 18,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MUL,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "11",
				},
			},
			want: &object.Int{
				Value: 77,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MOD,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "15",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
			},
			want: &object.Int{
				Value: 0,
			},
		},
		{
			in: &ast.UnaryExpr{
				Kind: lexer.MINUS,
				Value: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
			},
			want: &object.Int{
				Value: -3,
			},
		},
		{
			in: &ast.UnaryExpr{
				Kind: lexer.PLUS,
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3.0",
				},
			},
			want: &object.Float{
				Value: 3.0,
			},
		},
	}

	for _, c := range cases {
		got, err := evaluator.Eval(c.in)
		if err != nil {
			t.Fatalf("in(%v): %v", c.in, err)
		}

		if got.Type() != c.want.Type() {
			t.Errorf("got=%T, want=%T", got, c.want)
		}

		if got.String() != c.want.String() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
