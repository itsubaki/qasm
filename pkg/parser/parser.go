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
		case lexer.GATE, lexer.DEF:
			p.appendStmt(p.parseFunc())
		case lexer.IDENT:
			p.appendStmt(p.parse())
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

func (p *Parser) appendIncl(s ast.Stmt) {
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

func (p *Parser) parseIncl() ast.Stmt {
	c := p.next()
	p.expect(lexer.STRING)

	return &ast.InclStmt{
		Path: &ast.IdentExpr{
			Value: c.Literal,
		},
	}
}

func (p *Parser) parseIdentList() ast.ExprList {
	out := ast.ExprList{}
	out.Append(p.parseIdent())

	for p.cur.Token == lexer.COMMA {
		out.Append(p.parseIdent())
	}

	return out
}

func (p *Parser) parseIdent() ast.Expr {
	c := p.cur
	if p.cur.Token != lexer.IDENT {
		c = p.next()
	}
	p.expect(lexer.IDENT)

	x := &ast.IdentExpr{
		Value: c.Literal,
	}

	p.next()
	if p.cur.Token != lexer.LBRACKET {
		// q
		return x
	}

	v := p.next()
	lit := v.Literal
	if v.Token == lexer.MINUS {
		// q[-1]
		lit = fmt.Sprintf("%s%s", lit, p.next().Literal)
	}
	p.expect(lexer.INT)
	p.next()

	p.expect(lexer.RBRACKET)
	p.next()

	// q[0], q[-1]
	return &ast.IndexExpr{
		Name:  x,
		Value: lit,
	}
}

func (p *Parser) parseConst() ast.Stmt {
	p.expect(lexer.CONST)

	n := p.next()
	p.expect(lexer.IDENT)

	p.next()
	p.expect(lexer.EQUALS)

	v := p.next()
	p.expect(lexer.INT)

	// const N = 15
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
					Value: lexer.Tokens[kind],
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
			QArgs: p.parseIdentList(),
		},
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	p.expect(lexer.MEASURE)

	left := ast.MeasureExpr{
		QArgs: p.parseIdentList(),
	}

	if p.cur.Token != lexer.ARROW {
		// measure q
		return &ast.ExprStmt{
			X: &left,
		}
	}

	// measure q -> c
	right := p.parseIdent()
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
			QArgs: p.parseIdentList(),
		},
	}
}

func (p *Parser) parseApply() ast.Stmt {
	g := p.cur.Token

	x := ast.ApplyExpr{
		Kind: g,
	}

	p.next()
	if p.cur.Token == lexer.LPAREN {
		x.Params = p.parseIdentList()
		p.expect(lexer.RPAREN)
	}
	x.QArgs = p.parseIdentList()

	return &ast.ExprStmt{
		X: &x,
	}
}

func (p *Parser) parseFunc() ast.Stmt {
	kind := p.cur.Token

	ident := p.next()
	p.expect(lexer.IDENT)

	d := ast.FuncDecl{
		Kind: kind,
		Name: ident.Literal,
		Body: &ast.BlockStmt{},
	}

	p.next()
	if p.cur.Token == lexer.LPAREN {
		d.Params = p.parseIdentList()
		p.expect(lexer.RPAREN)
	}

	d.QArgs = p.parseIdentList()
	p.expect(lexer.LBRACE)
	p.next()

	for p.cur.Token != lexer.RBRACE {
		d.Body.List = append(d.Body.List, p.parseApply())
		p.next()
	}
	p.expect(lexer.RBRACE)

	return &ast.DeclStmt{
		Decl: &d,
	}
}

func (p *Parser) parse() ast.Stmt {
	ident := p.parseIdent()

	if p.cur.Token == lexer.EQUALS {
		// c = measure q
		p.next()
		p.expect(lexer.MEASURE)

		return &ast.AssignStmt{
			Left: ident,
			Right: &ast.MeasureExpr{
				QArgs: p.parseIdentList(),
			},
		}
	}

	// shor(a, N) r0, r1
	x := ast.CallExpr{
		Name: ident.String(),
	}

	if p.cur.Token == lexer.LPAREN {
		x.Params = p.parseIdentList()
		p.expect(lexer.RPAREN)
	}
	x.QArgs = p.parseIdentList()

	return &ast.ExprStmt{
		X: &x,
	}
}
