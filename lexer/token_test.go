package lexer_test

import (
	"testing"

	"github.com/itsubaki/qasm/lexer"
)

func TestIsModifiler(t *testing.T) {
	var cases = []struct {
		in   lexer.Token
		want bool
	}{
		{lexer.CTRL, true},
		{lexer.NEGCTRL, true},
		{lexer.INV, true},
		{lexer.POW, true},
		{lexer.IDENT, false},
	}

	for _, c := range cases {
		got := lexer.IsModifiler(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestIsBinaryOperator(t *testing.T) {
	var cases = []struct {
		in   lexer.Token
		want bool
	}{
		{lexer.PLUS, true},
		{lexer.MINUS, true},
		{lexer.MUL, true},
		{lexer.DIV, true},
		{lexer.MOD, true},
		{lexer.IDENT, false},
	}

	for _, c := range cases {
		got := lexer.IsBinaryOperator(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestIsBasicLit(t *testing.T) {
	var cases = []struct {
		in   lexer.Token
		want bool
	}{
		{lexer.IDENT, true},
		{lexer.STRING, true},
		{lexer.INT, true},
		{lexer.FLOAT, true},
		{lexer.CTRL, false},
	}

	for _, c := range cases {
		got := lexer.IsBasicLit(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestIsConst(t *testing.T) {
	var cases = []struct {
		in   lexer.Token
		want bool
	}{
		{lexer.PI, true},
		{lexer.TAU, true},
		{lexer.EULER, true},
		{lexer.IDENT, false},
	}

	for _, c := range cases {
		got := lexer.IsConst(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
