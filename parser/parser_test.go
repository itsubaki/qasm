package parser_test

import (
	"testing"

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
