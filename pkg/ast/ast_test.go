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
		Includes: []ast.Expr{
			&ast.IdentExpr{
				Kind:  lexer.STRING,
				Value: "\"stdgates.qasm\"",
			},
		},
		Statements: []ast.Stmt{
			&ast.LetStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Target: []ast.IdentExpr{
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
					Target: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
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
					LBRACKET: lexer.LBRACKET,
					RBRACKET: lexer.RBRACKET,
					Kind:     lexer.INT,
					Value:    "2",
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

func TestLetStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.LetStmt
		want string
	}{
		{
			ast.LetStmt{
				Kind: lexer.BIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
			},
			"bit c",
		},
		{
			ast.LetStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			"qubit q",
		},
		{
			ast.LetStmt{
				Kind: lexer.QUBIT,
				Index: &ast.IndexExpr{
					LBRACKET: lexer.LBRACKET,
					RBRACKET: lexer.RBRACKET,
					Kind:     lexer.INT,
					Value:    "2",
				},
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			"qubit[2] q",
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
				Target: []ast.IdentExpr{
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
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							LBRACKET: lexer.LBRACKET,
							RBRACKET: lexer.RBRACKET,
							Kind:     lexer.INT,
							Value:    "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							LBRACKET: lexer.LBRACKET,
							RBRACKET: lexer.RBRACKET,
							Kind:     lexer.INT,
							Value:    "1",
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
				Target: []ast.IdentExpr{
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
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							LBRACKET: lexer.LBRACKET,
							RBRACKET: lexer.RBRACKET,
							Kind:     lexer.INT,
							Value:    "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							LBRACKET: lexer.LBRACKET,
							RBRACKET: lexer.RBRACKET,
							Kind:     lexer.INT,
							Value:    "1",
						},
					},
				},
			},
			"x p[0], p[1]",
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
				Target: []ast.IdentExpr{
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
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							LBRACKET: lexer.LBRACKET,
							RBRACKET: lexer.RBRACKET,
							Kind:     lexer.INT,
							Value:    "0",
						},
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.IndexExpr{
							LBRACKET: lexer.LBRACKET,
							RBRACKET: lexer.RBRACKET,
							Kind:     lexer.INT,
							Value:    "1",
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
					Target: []ast.IdentExpr{
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
						LBRACKET: lexer.LBRACKET,
						RBRACKET: lexer.RBRACKET,
						Kind:     lexer.INT,
						Value:    "0",
					},
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Target: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
							Index: &ast.IndexExpr{
								LBRACKET: lexer.LBRACKET,
								RBRACKET: lexer.RBRACKET,
								Kind:     lexer.INT,
								Value:    "0",
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
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
