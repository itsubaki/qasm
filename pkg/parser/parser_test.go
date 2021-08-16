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
	fmt.Printf("%v\n", ast)

	if errs := p.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	// Output:
	// OPENQASM 3.0;
	// include "itsubaki/q.qasm";
	// qubit q[2];
	// bit c[2];
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

func TestParseIncludes(t *testing.T) {
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
		got := p.Parse().Include[0].String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestParseStmt(t *testing.T) {
	var cases = []struct {
		in   string
		want string
	}{
		{
			"qubit q",
			"qubit q",
		},
		{
			"qubit q[2]",
			"qubit q[2]",
		},
		{
			"reset q",
			"reset q",
		},
		{
			"x q",
			"x q",
		},
		{
			"x p, q",
			"x p, q",
		},
		{
			"x p[0], q[1]",
			"x p[0], q[1]",
		},
		{
			"measure q",
			"measure q",
		},
		{
			"measure q -> c",
			"measure q -> c",
		},
		{
			"c = measure q",
			"c = measure q",
		},
		{
			"print",
			"print",
		},
		{
			"print q",
			"print q",
		},
		{
			"print q, p",
			"print q, p",
		},
		{
			"print q, p[0]",
			"print q, p[0]",
		},
	}

	for _, c := range cases {
		p := parser.New(lexer.New(strings.NewReader(string(c.in))))
		got := p.Parse().Statement[0].String()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
