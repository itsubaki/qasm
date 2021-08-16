package lexer_test

import (
	"os"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/pkg/lexer"
)

func TestLexer(t *testing.T) {
	type Token struct {
		token   lexer.Token
		literal string
	}

	var cases = []struct {
		in   string
		want []Token
	}{
		{
			in: "../../testdata/gate.qasm",
			want: []Token{
				{lexer.GATE, "gate"},
				{lexer.IDENT, "bell"},
				{lexer.IDENT, "q0"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "q1"},
				{lexer.LBRACE, "{"},

				{lexer.H, "h"},
				{lexer.IDENT, "q0"},
				{lexer.SEMICOLON, ";"},

				{lexer.CX, "cx"},
				{lexer.IDENT, "q0"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "q1"},
				{lexer.SEMICOLON, ";"},

				{lexer.RBRACE, "}"},
			},
		},
		{
			in: "../../testdata/test_lexer.qasm",
			want: []Token{
				{lexer.OPENQASM, "OPENQASM"},
				{lexer.FLOAT, "3.0"},
				{lexer.SEMICOLON, ";"},

				{lexer.INCLUDE, "include"},
				{lexer.STRING, "\"itsubaki/q.qasm\""},
				{lexer.SEMICOLON, ";"},

				{lexer.QUBIT, "qubit"},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "2"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.BIT, "bit"},
				{lexer.IDENT, "c"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "2"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.RESET, "reset"},
				{lexer.IDENT, "q"},
				{lexer.SEMICOLON, ";"},

				{lexer.H, "h"},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.CX, "cx"},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.MEASURE, "measure"},
				{lexer.IDENT, "q"},
				{lexer.ARROW, "->"},
				{lexer.IDENT, "c"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENT, "c"},
				{lexer.EQUALS, "="},
				{lexer.MEASURE, "measure"},
				{lexer.IDENT, "q"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENT, "c"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.EQUALS, "="},
				{lexer.MEASURE, "measure"},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENT, "c"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.RBRACKET, "]"},
				{lexer.EQUALS, "="},
				{lexer.MEASURE, "measure"},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.EOF, ""},
			},
		},
	}

	for _, c := range cases {
		f, err := os.ReadFile(c.in)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}

		lex := lexer.New(strings.NewReader(string(f)))
		for _, w := range c.want {
			token, literal := lex.Tokenize()
			if token != w.token || literal != w.literal {
				t.Errorf("got=%v:%v, want=%v:%v", token, literal, w.token, w.literal)
			}
		}

		if len(lex.Errors()) != 0 {
			t.Errorf("errors=%v", lex.Errors())
		}
	}
}

func TestLexerTokenize(t *testing.T) {
	var cases = []struct {
		in   string
		want []lexer.Token
	}{
		{"1", []lexer.Token{lexer.INT}},
		{"-1", []lexer.Token{lexer.MINUS, lexer.INT}},
		{"100", []lexer.Token{lexer.INT, lexer.EOF}},
		{"10.0", []lexer.Token{lexer.FLOAT, lexer.EOF}},
		{"\"abc\"", []lexer.Token{lexer.STRING, lexer.EOF}},
		{"'abc'", []lexer.Token{lexer.STRING, lexer.EOF}},
		{"abc", []lexer.Token{lexer.IDENT, lexer.EOF}},
		{" \t\n", []lexer.Token{lexer.WHITESPACE}},
		{"\\", []lexer.Token{lexer.ILLEGAL}},
		{"\"a", []lexer.Token{lexer.STRING, lexer.EOF}},
		{"print", []lexer.Token{lexer.PRINT}},
		{"->", []lexer.Token{lexer.ARROW, lexer.EOF}},
	}

	for _, c := range cases {
		lex := lexer.New(strings.NewReader(c.in))
		for _, w := range c.want {
			if got, _ := lex.TokenizeIgnore(); got != w {
				t.Fail()
			}
		}

		if len(lex.Errors()) != 0 {
			t.Errorf("errors=%v", lex.Errors())
		}
	}
}
