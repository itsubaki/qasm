package parser_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/listener"
	"github.com/itsubaki/qasm/parser"
)

func TestParse(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text: `OPENQASM 3.0; qubit q;`,
		},
		{
			text:   `qubit[ q;`,
			errMsg: `1:8: mismatched input ';' expecting ']'`,
		},
	}

	for _, c := range cases {
		_, err := parser.Parse(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%q, want=%q", err.Error(), c.errMsg)
			}

			continue
		}
	}
}

func TestLexicalErrors(t *testing.T) {
	const text = "qubit q;\n$ U(pi, 0, pi) q;"

	cases := []struct {
		name  string
		parse func(string) error
	}{
		{name: "Parse", parse: func(text string) error { _, err := parser.Parse(text); return err }},
		{name: "StringTree", parse: func(text string) error { _, err := parser.StringTree(text); return err }},
	}

	for _, c := range cases {
		err := c.parse(text)

		var syntaxErr *listener.SyntaxError
		if !errors.As(err, &syntaxErr) {
			t.Errorf("%s: got=%v, want SyntaxError", c.name, err)
			continue
		}

		if syntaxErr.Line != 2 || syntaxErr.Column != 0 || !strings.Contains(syntaxErr.Message, "token recognition error") {
			t.Errorf("%s: unexpected lexical error: %v", c.name, err)
		}
	}
}

func TestStringTree(t *testing.T) {
	cases := []struct {
		text   string
		tree   string
		errMsg string
	}{
		{
			text: `OPENQASM 3.0; qubit q;`,
			tree: `(program (version OPENQASM 3.0 ;) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) <EOF>)`,
		},
		{
			text:   `qubit[ q;`,
			errMsg: `1:8: mismatched input ';' expecting ']'`,
		},
	}

	for _, c := range cases {
		tree, err := parser.StringTree(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%q, want=%q", err.Error(), c.errMsg)
			}

			continue
		}

		if tree != c.tree {
			t.Errorf("got=%q, want=%q", tree, c.tree)
		}
	}
}
