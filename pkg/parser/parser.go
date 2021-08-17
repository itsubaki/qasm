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
			Version:    "3.0",
			Includes:   make([]string, 0),
			Gates:      make([]ast.Expr, 0),
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
		case lexer.GATE:
			p.appendGate(p.parseGate())
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

func (p *Parser) appendIncl(s string) {
	p.qasm.Includes = append(p.qasm.Includes, s)
}

func (p *Parser) appendGate(e ast.Expr) {
	p.qasm.Gates = append(p.qasm.Gates, e)
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

func (p *Parser) parseInclude() string {
	c := p.next()
	p.expect(lexer.STRING)

	return c.Literal
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
	p.expect(lexer.LBRACKET)

	s := p.next()
	v := s.Literal
	if s.Token == lexer.MINUS {
		v = fmt.Sprintf("%s%s", v, p.next().Literal)
	}
	p.expect(lexer.INT)

	p.next()
	p.expect(lexer.RBRACKET)

	return &ast.IndexExpr{
		Kind:  lexer.INT,
		Value: v,
	}
}

func (p *Parser) parseConst() ast.Stmt {
	kind := p.cur.Token // lexer.CONST

	n := p.next()
	p.expect(lexer.IDENT)

	p.next()
	p.expect(lexer.EQUALS)

	v := p.next()
	p.expect(lexer.INT)

	return &ast.ConstStmt{
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

	ident := p.next()
	p.expect(lexer.IDENT)

	// qubit q
	p.next()
	if p.cur.Token != lexer.LBRACKET {
		return &ast.DeclStmt{
			Kind: kind,
			Name: &ast.IdentExpr{
				Kind:  ident.Token,
				Value: ident.Literal,
			},
		}
	}

	// qubit q[2]
	p.expect(lexer.LBRACKET)

	index := p.next()
	p.expect(lexer.INT)

	p.next()
	p.expect(lexer.RBRACKET)

	return &ast.DeclStmt{
		Kind: kind,
		Name: &ast.IdentExpr{
			Kind:  ident.Token,
			Value: ident.Literal,
			Index: &ast.IndexExpr{
				Kind:  lexer.INT,
				Value: index.Literal,
			},
		},
	}
}

func (p *Parser) parseReset() ast.Stmt {
	p.expect(lexer.RESET)

	return &ast.ResetStmt{
		Kind:  lexer.RESET,
		QArgs: p.parseIdentList(),
	}
}

func (p *Parser) parseApply() ast.Stmt {
	kind := p.cur.Token // lexer.X, lexer.Y, ..., lexer.CX, ...
	params := make([]ast.IdentExpr, 0)

	p.next()
	if p.cur.Token == lexer.LPAREN {
		params = p.parseIdentList()
		p.expect(lexer.RPAREN)
	}

	return &ast.ApplyStmt{
		Kind:   kind,
		Params: params,
		QArgs:  p.parseIdentList(),
	}
}

func (p *Parser) parseMeasure() ast.Stmt {
	p.expect(lexer.MEASURE)

	// measure q -> c
	left := ast.MeasureStmt{
		Kind:  lexer.MEASURE,
		QArgs: p.parseIdentList(),
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
		Kind:  lexer.PRINT,
		QArgs: p.parseIdentList(),
	}
}

func (p *Parser) parseGate() ast.Expr {
	name := p.next()
	p.expect(lexer.IDENT)
	p.next()

	params := make([]ast.IdentExpr, 0)
	if p.cur.Token == lexer.LPAREN {
		params = p.parseIdentList()
		p.expect(lexer.RPAREN)
		p.next()
	}
	p.expect(lexer.IDENT)

	args := p.parseIdentList()
	p.expect(lexer.LBRACE)

	stmts := make([]ast.Stmt, 0)
	for {
		p.next()
		if p.cur.Token == lexer.RBRACE {
			break
		}

		stmts = append(stmts, p.parseApply())
	}
	p.expect(lexer.RBRACE)

	return &ast.GateExpr{
		Kind:       lexer.GATE,
		Name:       name.Literal,
		Params:     params,
		QArgs:      args,
		Statements: stmts,
	}
}

func (p *Parser) parse() ast.Stmt {
	p.expect(lexer.IDENT)

	left := p.parseIdent()
	if p.cur.Token == lexer.IDENT {
		// bell q, p
		return &ast.CallStmt{
			Kind:   lexer.IDENT,
			Name:   left.String(),
			Params: make([]ast.IdentExpr, 0),
			QArgs:  p.parseIdentList(),
		}
	}

	if p.cur.Token == lexer.LPAREN {
		// shor(a, N) r0, r1
		return &ast.CallStmt{
			Kind:   lexer.IDENT,
			Name:   left.String(),
			Params: p.parseIdentList(),
			QArgs:  p.parseIdentList(),
		}
	}

	if p.cur.Token == lexer.EQUALS {
		// c = measure q
		p.next()
		p.expect(lexer.MEASURE)

		return &ast.AssignStmt{
			Kind: lexer.EQUALS,
			Left: &left,
			Right: &ast.MeasureStmt{
				Kind: lexer.MEASURE,
				QArgs: []ast.IdentExpr{
					p.parseIdent(),
				},
			},
		}
	}

	p.appendErr(fmt.Errorf("invalid token=%v, literal=%v", p.cur.Token, p.cur.Literal))
	return nil
}
