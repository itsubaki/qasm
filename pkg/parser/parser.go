package parser

import (
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Parser struct {
	l *lexer.Lexer
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
	}
}

func (p *Parser) Parse() *ast.OpenQASM {
	qasm := &ast.OpenQASM{
		Version:    3.0,
		Statements: make([]ast.Stmt, 0),
	}

	return qasm
}
