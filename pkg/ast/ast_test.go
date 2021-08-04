package ast_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExampleProgram_String() {
	p := &ast.Program{
		Statements: []ast.Stmt{
			&ast.LetStmt{
				Kind: lexer.QUBIT,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				Name: []ast.Ident{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Name: &ast.Ident{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
		},
	}

	fmt.Println(p)

	// Output:
	// qubit q;
	// reset q;
	// x q;
	// c = measure q;
}

func TestIdentString(t *testing.T) {
	var cases = []struct {
		in   ast.Ident
		want string
	}{
		{
			ast.Ident{
				Kind:  lexer.STRING,
				Value: "q",
			},
			"q",
		},
		{
			ast.Ident{
				Kind:  lexer.STRING,
				Value: "q",
				Index: &ast.Index{
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

func TestLetStmtString(t *testing.T) {
	var cases = []struct {
		in   ast.LetStmt
		want string
	}{
		{
			ast.LetStmt{
				Kind: lexer.BIT,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "c",
				},
			},
			"bit c",
		},
		{
			ast.LetStmt{
				Kind: lexer.QUBIT,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			"qubit q",
		},
		{
			ast.LetStmt{
				Kind: lexer.QUBIT,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
					Index: &ast.Index{
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
				Name: []ast.Ident{
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
				Name: []ast.Ident{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
					{
						Kind:  lexer.STRING,
						Value: "p",
						Index: &ast.Index{
							Kind:  lexer.INT,
							Value: "2",
						},
					},
				},
			},
			"reset q, p[2]",
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
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			"x q",
		},
		{
			ast.ApplyStmt{
				Kind: lexer.X,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "p",
					Index: &ast.Index{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
			},
			"x p[2]",
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
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			"measure q",
		},
		{
			ast.MeasureStmt{
				Kind: lexer.MEASURE,
				Name: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "p",
					Index: &ast.Index{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
			},
			"measure p[2]",
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
				Left: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Name: &ast.Ident{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			"c = measure q",
		},
		{
			ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.Ident{
					Kind:  lexer.STRING,
					Value: "c",
					Index: &ast.Index{
						Kind:  lexer.INT,
						Value: "2",
					},
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Name: &ast.Ident{
						Kind:  lexer.STRING,
						Value: "q",
						Index: &ast.Index{
							Kind:  lexer.INT,
							Value: "2",
						},
					},
				},
			},
			"c[2] = measure q[2]",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}