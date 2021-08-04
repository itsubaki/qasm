package ast

import (
	"bytes"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type Program struct {
	Path       string
	Statements []Stmt
}

func (p *Program) String() string {
	var buf bytes.Buffer

	for _, s := range p.Statements {
		buf.WriteString(s.String())
	}

	return buf.String()
}

type Node interface {
	Literal() string
	String() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Ident struct {
	Token lexer.Token
	Value string
}

func (i *Ident) exprNode() {}

func (i *Ident) Literal() string {
	return lexer.Tokens[i.Token]
}

func (i *Ident) String() string {
	return i.Value
}

type RegStmt struct {
	Token lexer.Token
	Name  *Ident
	Type  Expr
}

func (r *RegStmt) stmtNode() {}

func (r *RegStmt) Literal() string {
	return lexer.Tokens[r.Token]
}

func (r *RegStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(r.Literal())
	buf.WriteString(" ")

	buf.WriteString(r.Name.String())
	if r.Type != nil {
		buf.WriteString(r.Type.String())
	}

	buf.WriteString(";")

	return buf.String()
}

type Array struct {
	Lbrack lexer.Token
	Rbrack lexer.Token
	Token  lexer.Token
	Index  *Ident
}

func (a *Array) exprNode() {}

func (a *Array) Literal() string {
	return lexer.Tokens[a.Token]
}

func (a *Array) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[a.Lbrack])
	buf.WriteString(a.Index.Value)
	buf.WriteString(lexer.Tokens[a.Rbrack])

	return buf.String()
}

type ResetStmt struct {
	Token  lexer.Token
	Target *RegStmt
}

func (r *ResetStmt) stmtNode() {}

func (r *ResetStmt) Literal() string {
	return lexer.Tokens[r.Token]
}

func (r *ResetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(r.Literal())
	buf.WriteString(" ")

	buf.WriteString(r.Target.Name.Value)
	if r.Target.Type != nil {
		buf.WriteString(r.Target.Type.String())
	}

	buf.WriteString(";")

	return buf.String()
}
