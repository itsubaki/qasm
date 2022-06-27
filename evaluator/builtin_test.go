package evaluator_test

import (
	"testing"

	"github.com/itsubaki/qasm/evaluator"
	"github.com/itsubaki/qasm/lexer"
)

func TestBuildin(t *testing.T) {
	var cases = []struct {
		in   lexer.Token
		want bool
	}{
		{lexer.U, true},
		{lexer.X, true},
		{lexer.Y, true},
		{lexer.Z, true},
		{lexer.H, true},
		{lexer.T, true},
		{lexer.S, true},
		{lexer.GATE, false},
	}

	for _, c := range cases {
		_, got := evaluator.Builtin(c.in, []float64{1.0, 2.0, 3.0})
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
