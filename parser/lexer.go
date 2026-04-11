package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
)

// Lex lexes the input text and returns the list of tokens.
func Lex(text string) []antlr.Token {
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	return lexer.GetAllTokens()
}
