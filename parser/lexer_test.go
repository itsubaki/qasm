package parser_test

import (
	"testing"

	"github.com/itsubaki/qasm/parser"
)

func TestLex(t *testing.T) {
	cases := []struct {
		text   string
		tokens []string
		errMsg string
	}{
		{
			text: `OPENQASM 3.0; qubit q;`,
			tokens: []string{
				"OPENQASM", "3.0", ";",
				"qubit", "q", ";",
			},
		},
		{
			text:   `qubit[ q;`,
			errMsg: `1:8: mismatched input ';' expecting ']'`,
			tokens: []string{
				"qubit", "[", "q", ";",
			},
		},
	}

	for _, c := range cases {
		tokens := parser.Lex(c.text)
		if len(tokens) != len(c.tokens) {
			t.Errorf("got=%d, want=%d", len(tokens), len(c.tokens))
			continue
		}

		for i, token := range tokens {
			if token.GetText() != c.tokens[i] {
				t.Errorf("got=%q, want=%q", token.GetText(), c.tokens[i])
			}
		}
	}
}
