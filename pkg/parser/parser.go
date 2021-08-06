package parser

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Parser struct {
	l      *lexer.Lexer
	qasm   *ast.OpenQASM
	errors []error
	cur    lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
		qasm: &ast.OpenQASM{
			Version:    "3.0",
			Includes:   make([]ast.Expr, 0),
			Statements: make([]ast.Stmt, 0),
		},
		errors: make([]error, 0),
	}
}

func (p *Parser) Parse() *ast.OpenQASM {
	for {
		token, _ := p.l.Tokenize()
		switch token {
		case lexer.OPENQASM:
			p.qasm.Version = p.parseVersion()
		case lexer.INCLUDE:
			p.appendIncl(p.parseInclude())
		case lexer.QUBIT, lexer.BIT:
			p.appendStmt(p.parseLet(token))
		case lexer.RESET:
			p.appendStmt(p.parseReset())
		case lexer.MEASURE:
			p.appendStmt(p.parseMeasure())
		case lexer.PRINT:
			p.appendStmt(p.parsePrint())
		case lexer.X, lexer.Y, lexer.Z, lexer.H:
			p.appendStmt(p.parseApply(token))
		case lexer.EOF:
			return p.qasm
		}
	}
}

func (p *Parser) appendIncl(s ast.Expr) {
	p.qasm.Includes = append(p.qasm.Includes, s)
}

func (p *Parser) appendStmt(s ast.Stmt) {
	p.qasm.Statements = append(p.qasm.Statements, s)
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

func (p *Parser) parseInclude() ast.Expr {
	token, path := p.l.Tokenize()
	return &ast.IdentExpr{
		Kind:  token,
		Value: path,
	}
}

func (p *Parser) parseLet(kind lexer.Token) ast.Stmt {
	// qubit q
	token, value := p.l.Tokenize()
	if token == lexer.IDENT {
		return &ast.LetStmt{
			Kind: kind,
			Name: &ast.IdentExpr{
				Kind:  token,
				Value: value,
			},
		}
	}

	// qubit[2] q
	index, idxv := p.l.Tokenize()  // '2'
	rbtok, _ := p.l.Tokenize()     // ']'
	ident, value := p.l.Tokenize() // q

	return &ast.LetStmt{
		Kind: kind,
		Name: &ast.IdentExpr{
			Kind:  ident,
			Value: value,
		},
		Index: &ast.IndexExpr{
			LBRACKET: token,
			RBRACKET: rbtok,
			Kind:     index,
			Value:    idxv,
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
		Target: []ast.IdentExpr{
			{
				Kind:  token,
				Value: value,
			},
		},
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	token, value := p.l.Tokenize()
	return &ast.MeasureStmt{
		Kind: lexer.MEASURE,
		Target: []ast.IdentExpr{
			{
				Kind:  token,
				Value: value,
			},
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
