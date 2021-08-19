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
	cur    *Cursor
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
		qasm: &ast.OpenQASM{
			Version: "3.0",
			Incls:   make([]string, 0),
			Stmts:   make([]ast.Stmt, 0),
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
			p.appendIncl(p.parseIncl())
		case lexer.QUBIT, lexer.BIT:
			p.appendStmt(p.parseDecl())
		case lexer.CONST:
			p.appendStmt(p.parseConst())
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

	p.appendErr(fmt.Errorf("got=%v, want=%v", lexer.Tokens[p.cur.Token], lexer.Tokens[t]))
}

func (p *Parser) appendIncl(s string) {
	p.qasm.Incls = append(p.qasm.Incls, s)
}

func (p *Parser) appendStmt(s ast.Stmt) {
	p.qasm.Stmts = append(p.qasm.Stmts, s)
}

func (p *Parser) appendErr(e error) {
	p.errors = append(p.errors, e.Error())
}

func (p *Parser) parseVersion() string {
	c := p.next()
	p.expect(lexer.FLOAT)

	return c.Literal
}

func (p *Parser) parseIncl() string {
	c := p.next()
	p.expect(lexer.STRING)

	return c.Literal
}

func (p *Parser) parseExprList() ast.ExprList {
	out := ast.ExprList{}
	out.Append(p.parseExpr())

	for {
		if p.cur.Token != lexer.COMMA {
			break
		}

		out.Append(p.parseExpr())
	}

	return out
}

func (p *Parser) parseExpr() ast.Expr {
	c := p.cur
	if p.cur.Token != lexer.IDENT {
		c = p.next()
	}
	p.expect(lexer.IDENT)

	p.next()
	x := &ast.IdentExpr{
		Value: c.Literal,
	}

	if p.cur.Token != lexer.LBRACKET {
		return x
	}

	// TODO FIX
	return x
}

func (p *Parser) parseConst() ast.Stmt {
	p.expect(lexer.CONST)

	n := p.next()
	p.expect(lexer.IDENT)

	p.next()
	p.expect(lexer.EQUALS)

	v := p.next()
	p.expect(lexer.INT)

	return &ast.DeclStmt{
		Decl: &ast.GenConst{
			Name: &ast.IdentExpr{
				Value: n.Literal,
			},
			Value: v.Literal,
		},
	}
}

func (p *Parser) parseDecl() ast.Stmt {
	kind := p.cur.Token // lexer.QUBIT, lexer.BIT

	n := p.next()
	if p.cur.Token == lexer.IDENT {
		// qubit q
		p.expect(lexer.IDENT)
		return &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Kind: kind,
				Type: &ast.IdentExpr{
					Value: lexer.Tokens[kind],
				},
				Name: &ast.IdentExpr{
					Value: n.Literal,
				},
			},
		}
	}

	// qubit[2] q
	p.expect(lexer.LBRACKET)

	index := p.next()
	p.expect(lexer.INT)

	p.next()
	p.expect(lexer.RBRACKET)

	ident := p.next()
	p.expect(lexer.IDENT)

	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Kind: kind,
			Type: &ast.IndexExpr{
				Name: &ast.IdentExpr{
					Value: "qubit",
				},
				Value: index.Literal,
			},
			Name: &ast.IdentExpr{
				Value: ident.Literal,
			},
		},
	}
}

func (p *Parser) parseReset() ast.Stmt {
	p.expect(lexer.RESET)

	return &ast.ExprStmt{
		X: &ast.ResetExpr{
			QArgs: p.parseExprList(),
		},
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	p.expect(lexer.MEASURE)

	left := ast.MeasureExpr{
		QArgs: p.parseExprList(),
	}

	if p.cur.Token != lexer.ARROW {
		// measure q
		return &ast.ExprStmt{
			X: &left,
		}
	}

	// measure q -> c
	right := p.parseExpr()
	return &ast.ArrowStmt{
		Left:  &left,
		Right: right,
	}
}

func (p *Parser) parsePrint() ast.Stmt {
	p.expect(lexer.PRINT)

	c := p.next()
	if c.Token != lexer.IDENT {
		// print
		return &ast.ExprStmt{
			X: &ast.PrintExpr{},
		}
	}
	p.expect(lexer.IDENT)

	// print q, p;
	return &ast.ExprStmt{
		X: &ast.PrintExpr{
			QArgs: p.parseExprList(),
		},
	}
}

func (p *Parser) parseApply() ast.Stmt {
	g := p.cur.Token

	params := ast.ExprList{}
	p.next()
	if p.cur.Token == lexer.LPAREN {
		params = p.parseExprList()
		p.expect(lexer.RPAREN)
	}

	return &ast.ExprStmt{
		X: &ast.ApplyExpr{
			Kind:   g,
			Params: params,
			QArgs:  p.parseExprList(),
		},
	}
}
