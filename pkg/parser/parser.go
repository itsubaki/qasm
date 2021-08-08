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
	errors []string
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
		errors: make([]string, 0),
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *ast.OpenQASM {
	for {
		p.next()
		switch p.cur.Token {
		case lexer.OPENQASM:
			p.qasm.Version = p.parseVersion()
		case lexer.INCLUDE:
			p.appendIncl(p.parseInclude())
		case lexer.CONST:
			p.appendStmt(p.parseConstDecl())
		case lexer.QUBIT, lexer.BIT:
			p.appendStmt(p.parseDecl())
		case lexer.RESET:
			p.appendStmt(p.parseReset())
		case lexer.MEASURE:
			p.appendStmt(p.parseMeasure())
		case lexer.PRINT:
			p.appendStmt(p.parsePrint())
		case lexer.X, lexer.Y, lexer.Z:
			p.appendStmt(p.parseApply())
		case lexer.H, lexer.S, lexer.T:
			p.appendStmt(p.parseApply())
		case lexer.CX, lexer.CZ:
			p.appendStmt(p.parseApply())
		case lexer.CCX:
			p.appendStmt(p.parseApply())
		case lexer.SWAP, lexer.QFT, lexer.IQFT:
			p.appendStmt(p.parseApply())
		case lexer.CMODEXP2:
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

func (p *Parser) expect(t lexer.Token) {
	if p.cur.Token == t {
		return
	}

	p.appendErr(fmt.Errorf("%v not found", lexer.Tokens[t]))
}

func (p *Parser) appendIncl(s ast.Expr) {
	p.qasm.Includes = append(p.qasm.Includes, s)
}

func (p *Parser) appendStmt(s ast.Stmt) {
	p.qasm.Statements = append(p.qasm.Statements, s)
}

func (p *Parser) appendErr(e error) {
	p.errors = append(p.errors, e.Error())
}

func (p *Parser) parseVersion() string {
	p.next()
	p.expect(lexer.FLOAT)

	return p.cur.Literal
}

func (p *Parser) parseInclude() ast.Expr {
	p.next()
	p.expect(lexer.STRING)

	return &ast.IdentExpr{
		Kind:  p.cur.Token,
		Value: p.cur.Literal,
	}
}

func (p Parser) parseIdentList() []ast.IdentExpr {
	out := make([]ast.IdentExpr, 0)
	out = append(out, p.parseIdent())

	for {
		if p.cur.Token != lexer.COMMA {
			break
		}

		out = append(out, p.parseIdent())
	}

	return out
}

func (p *Parser) parseIdent() ast.IdentExpr {
	ident := p.next()
	p.expect(lexer.IDENT)

	expr := ast.IdentExpr{
		Kind:  ident.Token,
		Value: ident.Literal,
	}

	p.next()
	if p.cur.Token != lexer.LBRACKET {
		return expr
	}

	expr.Index = p.parseIndex()
	p.next()

	return expr
}

func (p *Parser) parseIndex() *ast.IndexExpr {
	lbracket := p.cur
	p.expect(lexer.LBRACKET)

	index := p.next()
	p.expect(lexer.INT)

	rbracket := p.next()
	p.expect(lexer.RBRACKET)

	return &ast.IndexExpr{
		LBRACKET: lbracket.Token,
		RBRACKET: rbracket.Token,
		Kind:     index.Token,
		Value:    index.Literal,
	}
}

func (p *Parser) parseConstDecl() ast.Stmt {
	kind := p.cur.Token // lexer.CONST

	n := p.next()
	p.expect(lexer.IDENT)

	p.next()
	p.expect(lexer.EQUALS)

	v := p.next()
	p.expect(lexer.INT)

	return &ast.DeclConstStmt{
		Kind: kind,
		Name: &ast.IdentExpr{
			Kind:  n.Token,
			Value: n.Literal,
		},
		Value: v.Literal,
	}
}

func (p *Parser) parseDecl() ast.Stmt {
	kind := p.cur.Token // lexer.QUBIT, lexer.BIT
	c := p.next()       // ident or lbracket

	// qubit q
	if p.cur.Token == lexer.IDENT {
		return &ast.DeclStmt{
			Kind: kind,
			Name: &ast.IdentExpr{
				Kind:  p.cur.Token,
				Value: p.cur.Literal,
			},
		}
	}

	// qubit[2] q
	p.expect(lexer.LBRACKET)

	index := p.next()
	p.expect(lexer.INT)

	brack := p.next()
	p.expect(lexer.RBRACKET)

	ident := p.next()
	p.expect(lexer.IDENT)

	return &ast.DeclStmt{
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
	p.expect(lexer.RESET)

	return &ast.ResetStmt{
		Kind:   p.cur.Token,
		Target: p.parseIdentList(),
	}
}

func (p *Parser) parseApply() ast.Stmt {
	return &ast.ApplyStmt{
		Kind:   p.cur.Token,
		Target: p.parseIdentList(),
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	p.expect(lexer.MEASURE)

	return &ast.MeasureStmt{
		Kind:   p.cur.Token,
		Target: p.parseIdentList(),
	}
}

func (p *Parser) parsePrint() ast.Stmt {
	p.expect(lexer.PRINT)

	return &ast.PrintStmt{
		Kind: lexer.PRINT,
	}
}
