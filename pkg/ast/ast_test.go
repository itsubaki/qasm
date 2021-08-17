package ast_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleOpenQASM_String() {
	p := &ast.OpenQASM{
		Version: "3.0",
		Includes: []string{
			"\"stdgates.qasm\"",
		},
		Statements: []ast.Stmt{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					QArgs: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
				},
			},
			&ast.ArrowStmt{
				Kind: lexer.ARROW,
				Left: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					QArgs: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
				},
				Right: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
			},
		},
	}

	fmt.Println(p)

	// Output:
	// OPENQASM 3.0;
	// include "stdgates.qasm";
	// qubit q;
	// reset q;
	// x q;
	// c = measure q;
	// measure q -> c;
}

func TestIdentExprString(t *testing.T) {
	var cases = []struct {
		in   ast.IdentExpr
		want string
	}{
		{
			ast.IdentExpr{
				Kind:  lexer.STRING,
				Value: "q",
			},
			"q",
		},
		{
			ast.IdentExpr{
				Kind:  lexer.STRING,
				Value: "q",
				Index: &ast.IndexExpr{
					Kind:  lexer.INT,
					Value: "2",
				},
			},
			"q[2]",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestGateExprString(t *testing.T) {
	var cases = []struct {
		in   ast.GateExpr
		want string
	}{
		{
			ast.GateExpr{
				Kind: lexer.GATE,
				Name: "bell",
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q0",
					},
					{
						Kind:  lexer.STRING,
						Value: "q1",
					},
				},
				Statements: []ast.Stmt{
					&ast.ApplyStmt{
						Kind: lexer.H,
						QArgs: []ast.IdentExpr{
							{
								Kind:  lexer.IDENT,
								Value: "q0",
							},
						},
					},
					&ast.ApplyStmt{
						Kind: lexer.CX,
						QArgs: []ast.IdentExpr{
							{
								Kind:  lexer.IDENT,
								Value: "q0",
							},
							{
								Kind:  lexer.IDENT,
								Value: "q1",
							},
						},
					},
				},
			},
			"gate bell q0, q1 { h q0; cx q0, q1; }",
		},
		{
			ast.GateExpr{
				Kind: lexer.GATE,
				Name: "shor",
				Params: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "a",
					},
					{
						Kind:  lexer.STRING,
						Value: "N",
					},
				},
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "r0",
					},
					{
						Kind:  lexer.STRING,
						Value: "r1",
					},
				},
				Statements: []ast.Stmt{
					&ast.ApplyStmt{
						Kind: lexer.X,
						QArgs: []ast.IdentExpr{
							{
								Kind:  lexer.IDENT,
								Value: "r1",
								Index: &ast.IndexExpr{
									Kind:  lexer.INT,
									Value: "-1",
								},
							},
						},
					},
					&ast.ApplyStmt{
						Kind: lexer.H,
						QArgs: []ast.IdentExpr{
							{
								Kind:  lexer.IDENT,
								Value: "r0",
							},
						},
					},
					&ast.ApplyStmt{
						Kind: lexer.CMODEXP2,
						QArgs: []ast.IdentExpr{
							{
								Kind:  lexer.IDENT,
								Value: "a",
							},
							{
								Kind:  lexer.IDENT,
								Value: "N",
							},
							{
								Kind:  lexer.IDENT,
								Value: "r0",
							},
							{
								Kind:  lexer.IDENT,
								Value: "r1",
							},
						},
					},
					&ast.ApplyStmt{
						Kind: lexer.IQFT,
						QArgs: []ast.IdentExpr{
							{
								Kind:  lexer.IDENT,
								Value: "r0",
							},
						},
					},
				},
			},
			"gate shor(a, N) r0, r1 { x r1[-1]; h r0; cmodexp2 a, N, r0, r1; iqft r0; }",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestCallStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.CallStmt
		want string
	}{
		{
			ast.CallStmt{
				Kind: lexer.IDENT,
				Name: "bell",
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q0",
					},
					{
						Kind:  lexer.STRING,
						Value: "q1",
					},
				},
			},
			"bell q0, q1",
		},
		{
			ast.CallStmt{
				Kind: lexer.IDENT,
				Name: "shor",
				Params: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "a",
					},
					{
						Kind:  lexer.STRING,
						Value: "N",
					},
				},
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "r0",
					},
					{
						Kind:  lexer.STRING,
						Value: "r1",
					},
				},
			},
			"shor(a, N) r0, r1",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestConstStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.Stmt
		want string
	}{
		{
			&ast.ConstStmt{
				Kind: lexer.CONST,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "N",
				},
				Value: "15",
			},
			"const N = 15",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDeclStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.Stmt
		want string
	}{
		{
			&ast.DeclStmt{
				Kind: lexer.BIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
			},
			"bit c",
		},
		{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			"qubit q",
		},
		{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
					Index: &ast.IndexExpr{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
			},
			"qubit q[2]",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestResetStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.ResetStmt
		want string
	}{
		{
			ast.ResetStmt{
				Kind: lexer.RESET,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			"reset q",
		},
		{
			ast.ResetStmt{
				Kind: lexer.RESET,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "1",
						},
					},
				},
			},
			"reset p[0], p[1]",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestApplyStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.ApplyStmt
		want string
	}{
		{
			ast.ApplyStmt{
				Kind: lexer.X,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			"x q",
		},
		{
			ast.ApplyStmt{
				Kind: lexer.X,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "1",
						},
					},
				},
			},
			"x p[0], p[1]",
		},
		{
			ast.ApplyStmt{
				Kind: lexer.CX,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "1",
						},
					},
				},
			},
			"cx p[0], p[1]",
		},
		{
			ast.ApplyStmt{
				Kind: lexer.CCX,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "1",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "2",
						},
					},
				},
			},
			"ccx p[0], p[1], p[2]",
		},
		{
			ast.ApplyStmt{
				Kind: lexer.CMODEXP2,
				Params: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "a",
					},
					{
						Kind:  lexer.IDENT,
						Value: "N",
					},
				},
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "r0",
					},
					{
						Kind:  lexer.IDENT,
						Value: "r1",
					},
				},
			},
			"cmodexp2(a, N) r0, r1",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMeasureStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.MeasureStmt
		want string
	}{
		{
			ast.MeasureStmt{
				Kind: lexer.MEASURE,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			"measure q",
		},
		{
			ast.MeasureStmt{
				Kind: lexer.MEASURE,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							Kind:  lexer.INT,
							Value: "1",
						},
					},
				},
			},
			"measure p[0], p[1]",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestAssignStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.AssignStmt
		want string
	}{
		{
			ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					QArgs: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
				},
			},
			"c = measure q",
		},
		{
			ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
					Index: &ast.IndexExpr{
						Kind:  lexer.INT,
						Value: "0",
					},
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					QArgs: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
							Index: &ast.IndexExpr{
								Kind:  lexer.INT,
								Value: "0",
							},
						},
					},
				},
			},
			"c[0] = measure q[0]",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestPrintStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.PrintStmt
		want string
	}{
		{
			ast.PrintStmt{
				Kind: lexer.PRINT,
			},
			"print",
		},
		{
			ast.PrintStmt{
				Kind: lexer.PRINT,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			"print q",
		},
		{
			ast.PrintStmt{
				Kind: lexer.PRINT,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
					},
				},
			},
			"print q, p",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
