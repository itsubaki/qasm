package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/listener"
)

// Parse parses the input text and returns the AST (abstract syntax tree) of the program.
func Parse(text string) (parser.IProgramContext, error) {
	listener := &listener.ErrorListener{}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(listener)

	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	p.RemoveErrorListeners()     // remove default error listeners
	p.AddErrorListener(listener) // add custom error listener

	program := p.Program()
	if len(listener.Errors) > 0 {
		return nil, listener.Errors[0]
	}

	return program, nil
}

// StringTree parses the input text and returns the string tree of the program.
func StringTree(text string) (string, error) {
	listener := &listener.ErrorListener{}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(listener)

	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	p.RemoveErrorListeners()     // remove default error listeners
	p.AddErrorListener(listener) // add custom error listener

	program := p.Program()
	if len(listener.Errors) > 0 {
		return "", listener.Errors[0]
	}

	return program.ToStringTree(nil, p), nil
}
