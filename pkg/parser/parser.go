package parser

import (
	"fmt"
	"strings"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Cursor struct {
	Token   lexer.Token
	Literal string
}

type Parser struct {
	l      *lexer.Lexer
	cur    *Cursor
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l:      l,
		errors: make([]string, 0),
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Parse() *ast.OpenQASM {
	var version string
	var incls []ast.Stmt
	var stmts []ast.Stmt

	for p.next().Token != lexer.EOF {
		switch p.cur.Token {
		case lexer.OPENQASM:
			version = p.parseVersion()
		case lexer.INCLUDE:
			incls = append(incls, p.parseIncl())
		default:
			stmts = append(stmts, p.parseStmt())
		}
	}

	return &ast.OpenQASM{
		Version: version,
		Incls:   incls,
		Stmts:   stmts,
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

	p.error(fmt.Errorf("got=%#v, want=%#v", p.cur, lexer.Tokens[t]))
}

func (p *Parser) expectSemi() {
	p.expect(lexer.SEMICOLON)
}

func (p *Parser) error(e error) {
	p.errors = append(p.errors, e.Error())
}

func (p *Parser) parseVersion() string {
	p.expect(lexer.OPENQASM)

	v := p.next()
	p.expect(lexer.FLOAT)

	p.next()
	p.expectSemi()

	return v.Literal
}

func (p *Parser) parseIncl() ast.Stmt {
	p.expect(lexer.INCLUDE)

	c := p.next()
	p.expect(lexer.STRING)

	p.next()
	p.expectSemi()

	return &ast.InclStmt{
		Path: &ast.IdentExpr{
			Value: c.Literal,
		},
	}
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur.Token {
	case lexer.QUBIT, lexer.BIT, lexer.CONST, lexer.GATE, lexer.DEF:
		return p.parseDeclStmt()
	case lexer.IDENT:
		switch p.cur.Literal {
		case "int", "float":
			return p.parseDeclStmt()
		default:
			return p.parseAssignOrCall()
		}
	case lexer.RESET:
		return p.parseResetStmt()
	case lexer.MEASURE:
		return p.parseMeasureStmt()
	case lexer.PRINT:
		return p.parsePrintStmt()
	case lexer.X, lexer.Y, lexer.Z,
		lexer.H, lexer.S, lexer.T,
		lexer.CX, lexer.CZ, lexer.CCX,
		lexer.SWAP, lexer.QFT, lexer.IQFT,
		lexer.CMODEXP2:
		return p.parseApplyStmt()
	}

	p.error(fmt.Errorf("invalid stmt token=%#v", p.cur))
	return nil
}

func (p *Parser) parseDeclStmt() ast.Stmt {
	decl := p.parseDecl()
	if p.cur.Token != lexer.RBRACE {
		p.next()
		p.expectSemi()
	}

	return &ast.DeclStmt{
		Decl: decl,
	}
}

func (p *Parser) parseDecl() ast.Decl {
	switch p.cur.Token {
	case lexer.QUBIT, lexer.BIT:
		return p.parseGenDecl()
	case lexer.IDENT:
		switch p.cur.Literal {
		case "int":
			p.cur.Token = lexer.INT
			return p.parseGenDecl()
		case "float":
			p.cur.Token = lexer.FLOAT
			return p.parseGenDecl()
		}
	case lexer.CONST:
		return p.parseGenConst()
	case lexer.GATE:
		return p.parseGate()
	case lexer.DEF:
		return p.parseFunc()
	}

	p.error(fmt.Errorf("invalid decl token=%#v", p.cur))
	return nil
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

	// q[0]
	p.next()
	return &ast.IndexExpr{
		Name:  x,
		Value: lit,
	}
}

func (p *Parser) parseGenConst() ast.Decl {
	p.expect(lexer.CONST)

	n := p.next()
	p.expect(lexer.IDENT)

	p.next()
	p.expect(lexer.EQUALS)

	v := p.next()
	p.expect(lexer.INT)

	// const N = 15
	return &ast.GenConst{
		Name: &ast.IdentExpr{
			Value: n.Literal,
		},
		Value: v.Literal,
	}
}

func (p *Parser) parseGenDecl() ast.Decl {
	kind := p.cur.Token // lexer.QUBIT, lexer.BIT, lexer.INT, lexer.FLOAT

	n := p.next()
	if p.cur.Token == lexer.IDENT {
		// qubit q
		return &ast.GenDecl{
			Kind: kind,
			Type: &ast.IdentExpr{
				Value: strings.ToLower(lexer.Tokens[kind]),
			},
			Name: &ast.IdentExpr{
				Value: n.Literal,
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

	return &ast.GenDecl{
		Kind: kind,
		Type: &ast.IndexExpr{
			Name: &ast.IdentExpr{
				Value: strings.ToLower(lexer.Tokens[kind]),
			},
			Value: index.Literal,
		},
		Name: &ast.IdentExpr{
			Value: ident.Literal,
		},
	}
}

func (p *Parser) parseGate() ast.Decl {
	p.expect(lexer.GATE)

	ident := p.next()
	p.expect(lexer.IDENT)

	d := ast.GateDecl{
		Name: ident.Literal,
		Body: &ast.BlockStmt{},
	}

	p.next()
	if p.cur.Token == lexer.LPAREN {
		d.Params = ast.ParenExpr{
			List: p.parseIdentList(),
		}
		p.expect(lexer.RPAREN)
	}

	d.QArgs = p.parseIdentList()
	p.expect(lexer.LBRACE)

	p.next()
	for p.cur.Token != lexer.RBRACE {
		d.Body.List = append(d.Body.List, p.parseApplyStmt())
		p.next()
	}
	p.expect(lexer.RBRACE)

	// gate bell q0, q1 { h q0; cx q0, q1; }
	return &d
}

func (p *Parser) parseFunc() ast.Decl {
	p.expect(lexer.DEF)

	ident := p.next()
	p.expect(lexer.IDENT)

	d := ast.FuncDecl{
		Name: ident.Literal,
		Body: &ast.BlockStmt{},
	}

	// TODO

	return &d
}

func (p *Parser) parseResetStmt() ast.Stmt {
	p.expect(lexer.RESET)

	qargs := p.parseIdentList()
	p.expectSemi()

	// reset q, p;
	return &ast.ExprStmt{
		X: &ast.ResetExpr{
			QArgs: qargs,
		},
	}
}

func (p *Parser) parseMeasureStmt() ast.Stmt {
	p.expect(lexer.MEASURE)

	left := ast.MeasureExpr{
		QArgs: p.parseIdentList(),
	}

	if p.cur.Token != lexer.ARROW {
		p.expectSemi()

		// measure q;
		return &ast.ExprStmt{
			X: &left,
		}
	}
	p.expect(lexer.ARROW)

	// measure q -> c;
	right := p.parseIdent()
	p.expectSemi()

	return &ast.ArrowStmt{
		Left:  &left,
		Right: right,
	}
}

func (p *Parser) parsePrintStmt() ast.Stmt {
	p.expect(lexer.PRINT)

	c := p.next()
	if c.Token != lexer.IDENT {
		// print;

		p.expectSemi()
		return &ast.ExprStmt{
			X: &ast.PrintExpr{},
		}
	}
	p.expect(lexer.IDENT)

	// print q, p;
	qargs := p.parseIdentList()
	p.expectSemi()

	return &ast.ExprStmt{
		X: &ast.PrintExpr{
			QArgs: qargs,
		},
	}
}

func (p *Parser) parseApplyStmt() ast.Stmt {
	kind := p.cur.Token

	x := ast.ApplyExpr{
		Kind: kind,
	}

	p.next()
	if p.cur.Token == lexer.LPAREN {
		x.Params = ast.ParenExpr{
			List: p.parseIdentList(),
		}
		p.expect(lexer.RPAREN)
	}

	x.QArgs = p.parseIdentList()
	p.expectSemi()

	// cx q[0], q[1];
	return &ast.ExprStmt{
		X: &x,
	}
}

func (p *Parser) parseAssignOrCall() ast.Stmt {
	ident := p.parseIdent()

	if p.cur.Token == lexer.EQUALS {
		// c = measure q;

		p.next()
		p.expect(lexer.MEASURE)

		qargs := p.parseIdentList()
		p.expectSemi()

		return &ast.AssignStmt{
			Left: ident,
			Right: &ast.MeasureExpr{
				QArgs: qargs,
			},
		}
	}

	// shor(a, N) r0, r1;
	x := ast.CallExpr{
		Name: ident.String(),
	}

	if p.cur.Token == lexer.LPAREN {
		x.Params = ast.ParenExpr{
			List: p.parseIdentList(),
		}
		p.expect(lexer.RPAREN)
	}

	x.QArgs = p.parseIdentList()
	p.expectSemi()

	return &ast.ExprStmt{
		X: &x,
	}
}
