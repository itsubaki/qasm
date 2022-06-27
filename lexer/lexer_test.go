package lexer_test

import (
	"strings"
	"testing"

	"github.com/itsubaki/qasm/lexer"
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
			in: `
OPENQASM 3.0;

gate h  q    { U(pi/2.0, 0.0, pi) q; }
gate x  q    { U(pi, 0.0, pi) q; }
gate cx q, p { ctrl @ x q, p; }

qubit[2] q;
bit[2]   c;
reset q;

h  q[0];
cx q[0], q[1];
measure q -> c;
`,
			want: []Token{
				{lexer.OPENQASM, "OPENQASM"},
				{lexer.FLOAT, "3.0"},
				{lexer.SEMICOLON, ";"},

				{lexer.GATE, "gate"},
				{lexer.IDENT, "h"},
				{lexer.IDENT, "q"},
				{lexer.LBRACE, "{"},
				{lexer.U, "U"},
				{lexer.LPAREN, "("},
				{lexer.PI, "pi"},
				{lexer.DIV, "/"},
				{lexer.FLOAT, "2.0"},
				{lexer.COMMA, ","},
				{lexer.FLOAT, "0.0"},
				{lexer.COMMA, ","},
				{lexer.PI, "pi"},
				{lexer.RPAREN, ")"},
				{lexer.IDENT, "q"},
				{lexer.SEMICOLON, ";"},
				{lexer.RBRACE, "}"},

				{lexer.GATE, "gate"},
				{lexer.IDENT, "x"},
				{lexer.IDENT, "q"},
				{lexer.LBRACE, "{"},
				{lexer.U, "U"},
				{lexer.LPAREN, "("},
				{lexer.PI, "pi"},
				{lexer.COMMA, ","},
				{lexer.FLOAT, "0.0"},
				{lexer.COMMA, ","},
				{lexer.PI, "pi"},
				{lexer.RPAREN, ")"},
				{lexer.IDENT, "q"},
				{lexer.SEMICOLON, ";"},
				{lexer.RBRACE, "}"},

				{lexer.GATE, "gate"},
				{lexer.IDENT, "cx"},
				{lexer.IDENT, "q"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "p"},
				{lexer.LBRACE, "{"},
				{lexer.CTRL, "ctrl"},
				{lexer.AT, "@"},
				{lexer.IDENT, "x"},
				{lexer.IDENT, "q"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "p"},
				{lexer.SEMICOLON, ";"},
				{lexer.RBRACE, "}"},

				{lexer.QUBIT, "qubit"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "2"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENT, "q"},
				{lexer.SEMICOLON, ";"},

				{lexer.BIT, "bit"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "2"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENT, "c"},
				{lexer.SEMICOLON, ";"},

				{lexer.RESET, "reset"},
				{lexer.IDENT, "q"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENT, "h"},
				{lexer.IDENT, "q"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "0"},
				{lexer.RBRACKET, "]"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENT, "cx"},
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
			},
		},
		{
			in: "gate bell q0, q1 { h q0; cx q0, q1; }",
			want: []Token{
				{lexer.GATE, "gate"},
				{lexer.IDENT, "bell"},
				{lexer.IDENT, "q0"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "q1"},
				{lexer.LBRACE, "{"},

				{lexer.IDENT, "h"},
				{lexer.IDENT, "q0"},
				{lexer.SEMICOLON, ";"},

				{lexer.IDENT, "cx"},
				{lexer.IDENT, "q0"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "q1"},
				{lexer.SEMICOLON, ";"},

				{lexer.RBRACE, "}"},
			},
		},
		{
			in: "def shor(int[32] a, int[32] N) qubit[n] r0, qubit[m] r1 -> bit[n] { h r0; CMODEXP2(a, N) r0, r1; IQFT r0; return measure r0; }",
			want: []Token{
				{lexer.DEF, "def"},
				{lexer.IDENT, "shor"},

				{lexer.LPAREN, "("},
				{lexer.IDENT, "int"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "32"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENT, "a"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "int"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "32"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENT, "N"},
				{lexer.RPAREN, ")"},

				{lexer.QUBIT, "qubit"},
				{lexer.LBRACKET, "["},
				{lexer.IDENT, "n"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENT, "r0"},
				{lexer.COMMA, ","},

				{lexer.QUBIT, "qubit"},
				{lexer.LBRACKET, "["},
				{lexer.IDENT, "m"},
				{lexer.RBRACKET, "]"},
				{lexer.IDENT, "r1"},

				{lexer.ARROW, "->"},

				{lexer.BIT, "bit"},
				{lexer.LBRACKET, "["},
				{lexer.IDENT, "n"},
				{lexer.RBRACKET, "]"},
				{lexer.LBRACE, "{"},

				{lexer.IDENT, "h"},
				{lexer.IDENT, "r0"},
				{lexer.SEMICOLON, ";"},

				{lexer.CMODEXP2, "CMODEXP2"},
				{lexer.LPAREN, "("},
				{lexer.IDENT, "a"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "N"},
				{lexer.RPAREN, ")"},
				{lexer.IDENT, "r0"},
				{lexer.COMMA, ","},
				{lexer.IDENT, "r1"},
				{lexer.SEMICOLON, ";"},

				{lexer.IQFT, "IQFT"},
				{lexer.IDENT, "r0"},
				{lexer.SEMICOLON, ";"},

				{lexer.RETURN, "return"},
				{lexer.MEASURE, "measure"},
				{lexer.IDENT, "r0"},
				{lexer.SEMICOLON, ";"},
			},
		},
	}

	for _, c := range cases {
		lex := lexer.New(strings.NewReader(c.in))
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
