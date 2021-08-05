package parser

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Parser struct {
	l      *lexer.Lexer
	errors []error
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l:      l,
		errors: make([]error, 0),
	}
}

func (p *Parser) Parse() *ast.OpenQASM {
	qasm := &ast.OpenQASM{
		Statements: make([]ast.Stmt, 0),
	}

	for {
		token, _ := p.l.Tokenize()
		switch token {
		case lexer.OPENQASM:
			qasm.Version = p.parseVersion()
		case lexer.QUBIT, lexer.BIT:
			stmt := p.parseLet(token)
			qasm.Statements = append(qasm.Statements, stmt)
		case lexer.RESET:
			stmt := p.parseReset()
			qasm.Statements = append(qasm.Statements, stmt)
		case lexer.MEASURE:
			stmt := p.parseMeasure()
			qasm.Statements = append(qasm.Statements, stmt)
		case lexer.PRINT:
			stmt := p.parsePrint()
			qasm.Statements = append(qasm.Statements, stmt)
		case lexer.X, lexer.Y, lexer.Z, lexer.H:
			stmt := p.parseApply(token)
			qasm.Statements = append(qasm.Statements, stmt)
		case lexer.EOF:
			return qasm
		}
	}
}

func (p *Parser) parseVersion() string {
	token, version := p.l.Tokenize()
	if token != lexer.FLOAT {
		msg := fmt.Errorf("invalid token=%v", version)
		p.errors = append(p.errors, msg)
		return ""
	}

	return version
}

func (p *Parser) parseLet(kind lexer.Token) ast.Stmt {
	token, value := p.l.Tokenize()
	if token != lexer.IDENT {
		msg := fmt.Errorf("invalid token=%v", value)
		p.errors = append(p.errors, msg)
	}

	return &ast.LetStmt{
		Kind: kind,
		Name: &ast.IdentExpr{
			Kind:  token,
			Value: value,
		},
	}
}

func (p *Parser) parseReset() ast.Stmt {
	token, value := p.l.Tokenize()

	return &ast.ResetStmt{
		Kind: lexer.RESET,
		Target: []ast.IdentExpr{
			{
				Kind:  token,
				Value: value,
			},
		},
	}
}

func (p *Parser) parseApply(kind lexer.Token) ast.Stmt {
	token, value := p.l.Tokenize()

	return &ast.ApplyStmt{
		Kind: kind,
		Target: &ast.IdentExpr{
			Kind:  token,
			Value: value,
		},
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	token, value := p.l.Tokenize()

	return &ast.MeasureStmt{
		Kind: lexer.MEASURE,
		Target: &ast.IdentExpr{
			Kind:  token,
			Value: value,
		},
	}
}

func (p *Parser) parsePrint() ast.Stmt {
	return &ast.PrintStmt{
		Kind: lexer.PRINT,
	}
}

func (p *Parser) Errors() []error {
	return p.errors
}
