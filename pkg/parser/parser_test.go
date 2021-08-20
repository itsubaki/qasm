package parser_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
)

func ExampleParser() {
	f, err := os.ReadFile("../../testdata/test_parser.qasm")
	if err != nil {
		fmt.Printf("read file: %v", err)
		return
	}

	p := parser.New(lexer.New(strings.NewReader(string(f))))
	ast := p.Parse()
	fmt.Println(ast)

	if errs := p.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	// Output:
	// OPENQASM 3.0;
	// include "itsubaki/q.qasm";
	// qubit[2] q;
	// bit[2] c;
	// reset q;
	// h q[0];
	// cx q[0], q[1];
	// measure q -> c;
	// c = measure q;
	// c[0] = measure q[0];
	// c[1] = measure q[1];
}

func TestParseVersion(t *testing.T) {
	var cases = []struct {
		in   string
		want string
	}{
		{
			"OPENQASM 3.0;",
			"3.0",
		},
	}

	for _, c := range cases {
		p := parser.New(lexer.New(strings.NewReader(string(c.in))))
		got := p.Parse().Version
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestParseIncl(t *testing.T) {
	var cases = []struct {
		in   string
		want string
	}{
		{
			"include \"gate01.qasm\";",
			"\"gate01.qasm\"",
		},
	}

	for _, c := range cases {
		p := parser.New(lexer.New(strings.NewReader(string(c.in))))
		got := p.Parse().Incls[0]
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestParseStmt(t *testing.T) {
	var cases = []struct {
		in string
	}{
		{"bit c;"},
		{"qubit q;"},
		{"qubit[2] q;"},
		{"reset q;"},
		{"reset q, p;"},
		{"reset q[0], p[0];"},
		{"measure q;"},
		{"measure q, p;"},
		{"measure q[0], q[1];"},
		{"print;"},
		{"print q;"},
		{"print q, p;"},
		{"print q, p[0];"},
		{"h q;"},
		{"h q, p;"},
		{"h q[0], p[0];"},
		{"h q[-1];"},
		{"gate bell q0, q1 { h q0; cx q0, q1; }"},
	}

	for _, c := range cases {
		p := parser.New(lexer.New(strings.NewReader(string(c.in))))
		got := p.Parse().Stmts[0].String()
		if got != c.in {
			t.Errorf("got=%v, want=%v", got, c.in)
		}
	}
}
