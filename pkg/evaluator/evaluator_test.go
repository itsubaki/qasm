package evaluator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
)

func Example_bell() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi/2.0) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

qubit[2] q;
reset    q;

h  q[0];
cx q[0], q[1];
`

	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		fmt.Printf("parse: %v\n", errs)
		return
	}

	e := evaluator.Default()
	if err := e.Eval(a); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Printf("print: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_bellGate() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi/2.0) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { h q; cx q, p; }

qubit[2] q;
reset    q;

bell q[0], q[1];
`

	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		fmt.Printf("parse: %v\n", errs)
		return
	}

	e := evaluator.Default()
	if err := e.Eval(a); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Printf("print: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_bellCtrl() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi/2.0) q; }
gate x q { U(pi, 0, pi) q; }

qubit[2] q;
reset    q;

h q[0];
ctrl @ x q[0], q[1];
`

	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		fmt.Printf("parse: %v\n", errs)
		return
	}

	e := evaluator.Default()
	if err := e.Eval(a); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Printf("print: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_negc() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate ncx  q0, q1     { negctrl @ x   q0, q1; }
gate cncx q0, q1, q2 { ctrl    @ ncx q0, q1, q2; }

qubit[3] q;
reset q;

x q[0];
cncx q[0], q[1], q[2];
`

	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		fmt.Printf("parse: %v\n", errs)
		return
	}

	e := evaluator.Default()
	if err := e.Eval(a); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Printf("print: %v\n", err)
		return
	}

	// Output:
	// [101][  5]( 1.0000 0.0000i): 1.0000
}

func Example_shor() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi/2.0) q; }

const N = 15;
const a = 7;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

x r1[-1];
h r0;
CMODEXP2(a, N) r0, r1;
IQFT r0;
`

	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		fmt.Printf("parse: %v\n", errs)
		return
	}

	e := evaluator.Default()
	if err := e.Eval(a); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	if err := e.Println(); err != nil {
		fmt.Printf("print: %v\n", err)
		return
	}

	// Output:
	// [000 0001][  0   1]( 0.2500 0.0000i): 0.0625
	// [000 0100][  0   4]( 0.2500 0.0000i): 0.0625
	// [000 0111][  0   7]( 0.2500 0.0000i): 0.0625
	// [000 1101][  0  13]( 0.2500 0.0000i): 0.0625
	// [010 0001][  2   1]( 0.2500 0.0000i): 0.0625
	// [010 0100][  2   4](-0.2500 0.0000i): 0.0625
	// [010 0111][  2   7]( 0.0000-0.2500i): 0.0625
	// [010 1101][  2  13]( 0.0000 0.2500i): 0.0625
	// [100 0001][  4   1]( 0.2500 0.0000i): 0.0625
	// [100 0100][  4   4]( 0.2500 0.0000i): 0.0625
	// [100 0111][  4   7](-0.2500 0.0000i): 0.0625
	// [100 1101][  4  13](-0.2500 0.0000i): 0.0625
	// [110 0001][  6   1]( 0.2500 0.0000i): 0.0625
	// [110 0100][  6   4](-0.2500 0.0000i): 0.0625
	// [110 0111][  6   7]( 0.0000 0.2500i): 0.0625
	// [110 1101][  6  13]( 0.0000-0.2500i): 0.0625
}

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
						&ast.IdentExpr{
							Value: "q",
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
	// X q;
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
					Name: "cx",
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
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Name: "cx",
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
