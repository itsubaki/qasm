package parser

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Cursor struct {
	Token   lexer.Token
	Literal string
}

type Parser struct {
	l      *lexer.Lexer
	qasm   *ast.OpenQASM
	errors []error
	cur    *Cursor
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
		p.next()
		switch p.cur.Token {
		case lexer.OPENQASM:
			p.qasm.Version = p.parseVersion()
		case lexer.INCLUDE:
			p.appendIncl(p.parseInclude())
		case lexer.QUBIT, lexer.BIT:
			p.appendStmt(p.parseLet())
		case lexer.RESET:
			p.appendStmt(p.parseReset())
		case lexer.MEASURE:
			p.appendStmt(p.parseMeasure())
		case lexer.PRINT:
			p.appendStmt(p.parsePrint())
		case lexer.X, lexer.Y, lexer.Z, lexer.H:
			p.appendStmt(p.parseApply())
		case lexer.EOF:
			return p.qasm
		}
	}
}

func (p *Parser) next() *Cursor {
	token, literal := p.l.Tokenize()
	p.cur = &Cursor{
		Token:   token,
		Literal: literal,
	}

	return p.cur

}

func (p *Parser) appendIncl(s ast.Expr) {
	p.qasm.Includes = append(p.qasm.Includes, s)
}

func (p *Parser) appendStmt(s ast.Stmt) {
	p.qasm.Statements = append(p.qasm.Statements, s)
}

func (p *Parser) appendErr(e error) {
	p.errors = append(p.errors, e)
}

func (p *Parser) parseVersion() string {
	c := p.next()
	if c.Token != lexer.FLOAT {
		p.appendErr(fmt.Errorf("invalid token=%v", c.Literal))
		return ""
	}

	return c.Literal
}

func (p *Parser) parseInclude() ast.Expr {
	c := p.next()
	return &ast.IdentExpr{
		Kind:  c.Token,
		Value: c.Literal,
	}
}

func (p *Parser) parseIdent() *ast.IdentExpr {
	return &ast.IdentExpr{}
}

func (p *Parser) parseIndex() *ast.IndexExpr {
	return &ast.IndexExpr{}
}

func (p *Parser) parseLet() ast.Stmt {
	kind := p.cur.Token

	// qubit q
	c := p.next()
	if c.Token == lexer.IDENT {
		return &ast.LetStmt{
			Kind: kind,
			Name: &ast.IdentExpr{
				Kind:  c.Token,
				Value: c.Literal,
			},
		}
	}

	// qubit[2] q
	// TODO check token
	index := p.next() // '2'
	brack := p.next() // ']'
	ident := p.next() // q

	return &ast.LetStmt{
		Kind: kind,
		Name: &ast.IdentExpr{
			Kind:  ident.Token,
			Value: ident.Literal,
		},
		Index: &ast.IndexExpr{
			LBRACKET: c.Token,
			RBRACKET: brack.Token,
			Kind:     index.Token,
			Value:    index.Literal,
		},
	}
}

func (p *Parser) parseReset() ast.Stmt {
	c := p.next()

	return &ast.ResetStmt{
		Kind: lexer.RESET,
		Target: []ast.IdentExpr{
			{
				Kind:  c.Token,
				Value: c.Literal,
			},
		},
	}
}

func (p *Parser) parseApply() ast.Stmt {
	kind := p.cur.Token
	c := p.next()

	return &ast.ApplyStmt{
		Kind: kind,
		Target: []ast.IdentExpr{
			{
				Kind:  c.Token,
				Value: c.Literal,
			},
		},
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	c := p.next()

	return &ast.MeasureStmt{
		Kind: lexer.MEASURE,
		Target: []ast.IdentExpr{
			{
				Kind:  c.Token,
				Value: c.Literal,
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
