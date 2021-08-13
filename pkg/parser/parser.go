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
		case lexer.QUBIT, lexer.BIT, lexer.CONST:
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
		case lexer.IDENT:
			p.appendStmt(p.parseAssign())
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
	c := p.next()
	p.expect(lexer.FLOAT)

	return c.Literal
}

func (p *Parser) parseInclude() ast.Expr {
	c := p.next()
	p.expect(lexer.STRING)

	return &ast.IdentExpr{
		Kind:  c.Token,
		Value: c.Literal,
	}
}

func (p *Parser) parseIdentList() []ast.IdentExpr {
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
	c := p.cur
	if p.cur.Token != lexer.IDENT {
		c = p.next()
	}
	p.expect(lexer.IDENT)

	ident := ast.IdentExpr{
		Kind:  c.Token,
		Value: c.Literal,
	}

	p.next()
	if p.cur.Token != lexer.LBRACKET {
		return ident
	}

	ident.Index = p.parseIndex()
	p.next()

	return ident
}

func (p *Parser) parseIndex() *ast.IndexExpr {
	lbrack := p.cur
	p.expect(lexer.LBRACKET)

	s := p.next()
	v := s.Literal
	if s.Token == lexer.MINUS {
		v = fmt.Sprintf("%s%s", v, p.next().Literal)
	}
	p.expect(lexer.INT)

	rbrack := p.next()
	p.expect(lexer.RBRACKET)

	return &ast.IndexExpr{
		LBRACKET: lbrack.Token,
		RBRACKET: rbrack.Token,
		Kind:     lexer.INT,
		Value:    v,
	}
}

func (p *Parser) parseDecl() ast.Stmt {
	kind := p.cur.Token // lexer.QUBIT, lexer.BIT, lexer.CONST

	if kind == lexer.CONST {
		n := p.next()
		p.expect(lexer.IDENT)

		p.next()
		p.expect(lexer.EQUALS)

		v := p.next()
		p.expect(lexer.INT)

		return &ast.DeclStmt{
			Kind: kind,
			Name: &ast.IdentExpr{
				Kind:  n.Token,
				Value: n.Literal,
			},
			Value: v.Literal,
		}
	}

	c := p.next() // ident or lbracket

	// qubit q
	if c.Token == lexer.IDENT {
		return &ast.DeclStmt{
			Kind: kind,
			Name: &ast.IdentExpr{
				Kind:  c.Token,
				Value: c.Literal,
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
		Kind:   lexer.RESET,
		Target: p.parseIdentList(),
	}
}

func (p *Parser) parseApply() ast.Stmt {
	kind := p.cur.Token // lexer.X, lexer.Y, ..., lexer.CX, ...
	return &ast.ApplyStmt{
		Kind:   kind,
		Target: p.parseIdentList(),
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	p.expect(lexer.MEASURE)

	// measure q -> c
	left := ast.MeasureStmt{
		Kind:   lexer.MEASURE,
		Target: p.parseIdentList(),
	}

	if p.cur.Token != lexer.ARROW {
		return &left
	}

	right := p.parseIdent()
	return &ast.ArrowStmt{
		Kind:  lexer.ARROW,
		Left:  &left,
		Right: &right,
	}
}

func (p *Parser) parseAssign() ast.Stmt {
	p.expect(lexer.IDENT)

	// c = measure q
	left := p.parseIdent()
	p.expect(lexer.EQUALS)

	p.next()
	p.expect(lexer.MEASURE)

	right := ast.MeasureStmt{
		Kind: lexer.MEASURE,
		Target: []ast.IdentExpr{
			p.parseIdent(),
		},
	}

	return &ast.AssignStmt{
		Kind:  lexer.EQUALS,
		Left:  &left,
		Right: &right,
	}
}

func (p *Parser) parsePrint() ast.Stmt {
	p.expect(lexer.PRINT)

	c := p.next()
	if c.Token != lexer.IDENT {
		return &ast.PrintStmt{
			Kind: lexer.PRINT,
		}
	}
	p.expect(lexer.IDENT)

	return &ast.PrintStmt{
		Kind:   lexer.PRINT,
		Target: p.parseIdentList(),
	}
}
