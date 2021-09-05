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
	peek   *Cursor
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
	var version ast.Stmt
	var stmts []ast.Stmt

	p.next()
	for p.next().Token != lexer.EOF {
		switch p.cur.Token {
		case lexer.OPENQASM:
			version = p.parseVersion()
		default:
			stmts = append(stmts, p.parseStmt())
		}
	}

	return &ast.OpenQASM{
		Version: version,
		Stmts:   stmts,
	}
}

func (p *Parser) next() *Cursor {
	token, literal := p.l.Tokenize()

	p.cur = p.peek
	p.peek = &Cursor{
		Token:   token,
		Literal: literal,
	}

	return p.cur
}

func (p *Parser) expect(t lexer.Token) {
	if p.cur.Token == t {
		return
	}

	p.error(fmt.Errorf("got={Token:%v, Literal: %v}, want={Token:%v, Literal: %v}", p.cur.Token, p.cur.Literal, t, lexer.Tokens[t]))
}

func (p *Parser) expectSemi() {
	p.expect(lexer.SEMICOLON)
}

func (p *Parser) error(e error) {
	p.errors = append(p.errors, e.Error())
}

func (p *Parser) parseVersion() ast.Stmt {
	p.expect(lexer.OPENQASM)

	v := p.next()
	p.expect(lexer.FLOAT)

	p.next()
	p.expectSemi()

	return &ast.DeclStmt{
		Decl: &ast.VersionDecl{
			Value: &ast.BasicLit{
				Kind:  lexer.FLOAT,
				Value: v.Literal,
			},
		},
	}
}

func (p *Parser) parseIncl() ast.Stmt {
	p.expect(lexer.INCLUDE)

	c := p.next()
	p.expect(lexer.STRING)

	p.next()
	p.expectSemi()

	return &ast.InclStmt{
		Path: ast.BasicLit{
			Kind:  lexer.STRING,
			Value: c.Literal,
		},
	}
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur.Token {
	case lexer.INCLUDE:
		return p.parseIncl()
	case lexer.QUBIT, lexer.BIT, lexer.CONST,
		lexer.GATE, lexer.DEF:
		return p.parseDeclStmt()
	case lexer.IDENT:
		switch p.cur.Literal {
		case "int", "float":
			return p.parseDeclStmt()
		}
		return p.parseAssignOrCall()
	case lexer.MEASURE:
		return p.parseMeasureStmt()
	case lexer.RESET:
		return p.parseResetStmt()
	case lexer.PRINT:
		return p.parsePrintStmt()
	case lexer.RETURN:
		return p.parseReturnStmt()
	case lexer.X, lexer.Y, lexer.Z, lexer.U,
		lexer.H, lexer.S, lexer.T,
		lexer.CX, lexer.CZ, lexer.CCX,
		lexer.SWAP, lexer.QFT, lexer.IQFT, lexer.CMODEXP2,
		lexer.CTRL, lexer.NEGCTRL, lexer.INV, lexer.POW:
		return p.parseApplyStmt()
	}

	p.error(fmt.Errorf("invalid stmt token=%#v", p.cur))
	return &ast.BadStmt{}
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

func (p *Parser) parseDeclList() ast.DeclList {
	list := ast.DeclList{}
	list.Append(p.parseDecl())

	for p.next().Token == lexer.COMMA {
		p.next() // skip COMMA token

		list.Append(p.parseDecl())
	}

	return list
}

func (p *Parser) parseDecl() ast.Decl {
	switch p.cur.Token {
	case lexer.QUBIT, lexer.BIT:
		return p.parseGenDecl()
	case lexer.IDENT:
		switch p.cur.Literal {
		case "int":
			p.cur.Token = lexer.INT
		case "float":
			p.cur.Token = lexer.FLOAT
		}
		return p.parseGenDecl()
	case lexer.CONST:
		return p.parseGenConst()
	case lexer.GATE:
		return p.parseGate()
	case lexer.DEF:
		return p.parseFunc()
	}

	p.error(fmt.Errorf("invalid decl token=%#v", p.cur))
	return &ast.BadDecl{}
}

func (p *Parser) parseIdentList() ast.ExprList {
	list := ast.ExprList{}
	list.Append(p.parseIdent())

	for p.next().Token == lexer.COMMA {
		p.next() // skip COMMNA token

		list.Append(p.parseIdent())
	}

	return list
}

func (p *Parser) parseInfix() ast.Expr {
	lhs := p.cur
	ope := p.next()
	rhs := p.next()

	return &ast.InfixExpr{
		Kind: ope.Token,
		Left: &ast.BasicLit{
			Kind:  lhs.Token,
			Value: lhs.Literal,
		},
		Right: &ast.BasicLit{
			Kind:  rhs.Token,
			Value: rhs.Literal,
		},
	}
}

func (p *Parser) parseIdent() ast.Expr {
	c := p.cur
	if !lexer.IsBasicLit(c.Token) {
		c = p.next()
	}

	if c.Token != lexer.IDENT && lexer.IsBasicLit(c.Token) {
		if lexer.IsBinaryOperator(p.peek.Token) {
			// pi / 2
			return p.parseInfix()
		}

		// pi, 1.23
		return &ast.BasicLit{
			Kind:  c.Token,
			Value: c.Literal,
		}
	}
	p.expect(lexer.IDENT)

	x := ast.IdentExpr{
		Value: c.Literal,
	}

	if p.peek.Token != lexer.LBRACKET {
		// q
		return &x
	}

	p.next()
	p.expect(lexer.LBRACKET)

	v := p.next().Literal
	if p.cur.Token == lexer.MINUS {
		// q[-1]
		v = fmt.Sprintf("%s%s", lexer.Tokens[lexer.MINUS], p.next().Literal)
	}
	p.expect(lexer.INT)

	p.next()
	p.expect(lexer.RBRACKET)

	// q[0]
	return &ast.IndexExpr{
		Name:  x,
		Value: v,
	}
}

func (p *Parser) parseGenConst() ast.Decl {
	p.expect(lexer.CONST)

	n := p.next()
	p.expect(lexer.IDENT)

	p.next()
	p.expect(lexer.EQUALS)

	v := p.next()
	if lexer.IsBinaryOperator(p.peek.Token) {
		// const N = pi * 2
		return &ast.GenConst{
			Name: ast.IdentExpr{
				Value: n.Literal,
			},
			Value: p.parseInfix(),
		}
	}

	// const N = 15
	return &ast.GenConst{
		Name: ast.IdentExpr{
			Value: n.Literal,
		},
		Value: &ast.BasicLit{
			Kind:  v.Token,
			Value: v.Literal,
		},
	}
}

func (p *Parser) parseGenDecl() ast.Decl {
	kind := p.cur.Token // lexer.QUBIT, lexer.BIT, lexer.INT, lexer.FLOAT

	if p.next().Token == lexer.IDENT {
		// qubit q
		return &ast.GenDecl{
			Kind: kind,
			Type: &ast.IdentExpr{
				Value: strings.ToLower(lexer.Tokens[kind]),
			},
			Name: ast.IdentExpr{
				Value: p.cur.Literal,
			},
		}
	}

	// qubit[2] q
	p.expect(lexer.LBRACKET)

	index := p.next()
	// p.expect(lexer.INT), or lexer.IDENT

	p.next()
	p.expect(lexer.RBRACKET)

	v := p.next()
	p.expect(lexer.IDENT)

	return &ast.GenDecl{
		Kind: kind,
		Type: &ast.IndexExpr{
			Name: ast.IdentExpr{
				Value: strings.ToLower(lexer.Tokens[kind]),
			},
			Value: index.Literal,
		},
		Name: ast.IdentExpr{
			Value: v.Literal,
		},
	}
}

func (p *Parser) parseGate() ast.Decl {
	p.expect(lexer.GATE)

	ident := p.next()
	p.expect(lexer.IDENT)

	decl := ast.GateDecl{
		Name: ident.Literal,
		Body: ast.BlockStmt{},
	}

	if p.next().Token == lexer.LPAREN {
		decl.Params = ast.ParenExpr{
			List: p.parseIdentList(),
		}
		p.expect(lexer.RPAREN)
	}

	decl.QArgs = p.parseIdentList()
	p.expect(lexer.LBRACE)

	for p.next().Token != lexer.RBRACE {
		decl.Body.Append(p.parseApplyStmt())
	}
	p.expect(lexer.RBRACE)

	// gate bell q0, q1 { h q0; cx q0, q1; }
	return &decl
}

func (p *Parser) parseFunc() ast.Decl {
	p.expect(lexer.DEF)

	ident := p.next()
	p.expect(lexer.IDENT)

	// def shor
	decl := ast.FuncDecl{
		Name: ident.Literal,
		Body: ast.BlockStmt{},
	}

	if p.next().Token == lexer.LPAREN {
		p.next() // skip LPAREN

		// (int[32] a, int[32] N)
		decl.Params = ast.ParenDecl{
			List: p.parseDeclList(),
		}

		p.expect(lexer.RPAREN)
		p.next()
	}

	// qubit[n] q0, qubit[m] q1
	decl.QArgs = p.parseDeclList()
	p.expect(lexer.ARROW) // ->

	// bit[n]
	bit := p.next()
	p.expect(lexer.BIT)

	p.next()
	p.expect(lexer.LBRACKET)

	val := p.next()

	p.next()
	p.expect(lexer.RBRACKET)

	decl.Result = &ast.IndexExpr{
		Name: ast.IdentExpr{
			Value: bit.Literal,
		},
		Value: val.Literal,
	}

	// { h q0; cx q0, q1; return measure q0, q1; }
	p.next()
	p.expect(lexer.LBRACE)

	for p.next().Token != lexer.RBRACE {
		decl.Body.Append(p.parseStmt())
	}
	p.expect(lexer.RBRACE)

	return &decl
}

func (p *Parser) parseCall(name string) ast.Expr {
	x := ast.CallExpr{
		Name: name,
	}

	if p.cur.Token == lexer.LPAREN {
		x.Params = ast.ParenExpr{
			List: p.parseIdentList(),
		}
		p.expect(lexer.RPAREN)
	}

	x.QArgs = p.parseIdentList()
	return &x
}

func (p *Parser) parseMeasure() ast.Expr {
	p.expect(lexer.MEASURE)

	return &ast.MeasureExpr{
		QArgs: p.parseIdentList(),
	}
}

func (p *Parser) parseMeasureStmt() ast.Stmt {
	p.expect(lexer.MEASURE)

	x := p.parseMeasure()
	if p.cur.Token != lexer.ARROW {
		p.expectSemi()

		// measure q;
		return &ast.ExprStmt{
			X: x,
		}
	}
	p.expect(lexer.ARROW)

	// measure q -> c;
	right := p.parseIdent()
	p.next()
	p.expectSemi()

	return &ast.ArrowStmt{
		Left:  x,
		Right: right,
	}
}

func (p *Parser) parseReturnStmt() ast.Stmt {
	p.expect(lexer.RETURN)

	if p.next().Token == lexer.MEASURE {
		x := p.parseMeasure()
		p.expectSemi()

		return &ast.ReturnStmt{
			Result: x,
		}
	}

	return &ast.ReturnStmt{
		Result: nil,
	}
}

func (p *Parser) parseResetStmt() ast.Stmt {
	p.expect(lexer.RESET)

	qargs := p.parseIdentList()
	p.expectSemi()

	// reset q, p;
	return &ast.ResetStmt{
		QArgs: qargs,
	}
}

func (p *Parser) parsePrintStmt() ast.Stmt {
	p.expect(lexer.PRINT)

	if p.next().Token != lexer.IDENT {
		p.expectSemi()

		// print;
		return &ast.PrintStmt{}
	}
	p.expect(lexer.IDENT)

	// print q, p;
	qargs := p.parseIdentList()
	p.expectSemi()

	return &ast.PrintStmt{
		QArgs: qargs,
	}
}

func (p *Parser) parseApplyStmt() ast.Stmt {
	mod := make([]ast.Modifier, 0)
	for lexer.IsModifiler(p.cur.Token) {
		m := ast.Modifier{
			Kind: p.cur.Token,
		}

		if p.next().Token == lexer.LPAREN {
			v := p.next()
			p.expect(lexer.INT)

			m.Index = ast.ParenExpr{
				List: ast.ExprList{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  v.Token,
							Value: v.Literal,
						},
					},
				},
			}

			p.next()
			p.expect(lexer.RPAREN)
			p.next()
		}
		mod = append(mod, m)

		p.expect(lexer.AT)
		p.next()
	}

	x := ast.ApplyStmt{
		Kind:     p.cur.Token,
		Name:     p.cur.Literal,
		Modifier: mod,
	}

	if p.next().Token == lexer.LPAREN {
		x.Params = ast.ParenExpr{
			List: p.parseIdentList(),
		}
		p.expect(lexer.RPAREN)
	}

	x.QArgs = p.parseIdentList()
	p.expectSemi()

	// cx q[0], q[1];
	return &x
}

func (p *Parser) parseAssignOrCall() ast.Stmt {
	c := p.parseIdent()

	if p.next().Token != lexer.EQUALS {
		// bell r0, r1;
		x := p.parseCall(c.String())
		p.expectSemi()

		return &ast.ExprStmt{
			X: x,
		}
	}
	p.expect(lexer.EQUALS)

	switch p.next().Token {
	case lexer.IDENT:
		p.expect(lexer.IDENT)
		n := p.cur.Literal
		p.next()

		// c = shor(a, N) r0, r1;
		x := p.parseCall(n)
		p.expectSemi()

		return &ast.AssignStmt{
			Left:  c,
			Right: x,
		}
	case lexer.MEASURE:
		p.expect(lexer.MEASURE)

		// c = measure q;
		m := p.parseMeasure()
		p.expectSemi()

		return &ast.AssignStmt{
			Left:  c,
			Right: m,
		}
	}

	p.error(fmt.Errorf("invalid assign token=%#v", p.cur))
	return &ast.BadStmt{}
}
