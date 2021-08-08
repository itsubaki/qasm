package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type OpenQASM struct {
	Version    string
	Includes   []Expr
	Statements []Stmt
}

func (p *OpenQASM) String() string {
	var buf bytes.Buffer

	version := fmt.Sprintf("OPENQASM %v;\n", p.Version)
	buf.WriteString(version)

	for _, i := range p.Includes {
		str := fmt.Sprintf("include %s;\n", i)
		buf.WriteString(str)
	}

	for _, s := range p.Statements {
		str := fmt.Sprintf("%s;\n", s.String())
		buf.WriteString(str)
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

type IdentExpr struct {
	Kind  lexer.Token // lexer.STRING, lexer.INT, lexer.FLOAT
	Value string
	Index *IndexExpr
}

func (i *IdentExpr) exprNode() {}

func (i *IdentExpr) Literal() string {
	return lexer.Tokens[i.Kind]
}

func (i *IdentExpr) String() string {
	if i.Index == nil {
		return i.Value
	}

	return fmt.Sprintf("%s%s", i.Value, i.Index.String())
}

type IndexExpr struct {
	LBRACKET lexer.Token
	RBRACKET lexer.Token
	Kind     lexer.Token // lexer.INT
	Value    string
}

func (i *IndexExpr) exprNode() {}

func (i *IndexExpr) Literal() string {
	return lexer.Tokens[i.Kind]
}

func (i *IndexExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[i.LBRACKET])
	buf.WriteString(i.Value)
	buf.WriteString(lexer.Tokens[i.RBRACKET])

	return buf.String()
}

func (i *IndexExpr) Int() int {
	v, err := strconv.Atoi(i.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type DeclStmt struct {
	Kind  lexer.Token // lexer.QUBIT, lexer.BIT
	Name  *IdentExpr
	Index *IndexExpr
}

func (s *DeclStmt) stmtNode() {}

func (s *DeclStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *DeclStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if s.Index != nil {
		buf.WriteString(s.Index.String())
	}
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())

	return buf.String()
}

type DeclConstStmt struct {
	Kind  lexer.Token // lexer.CONST
	Name  *IdentExpr
	Value string
}

func (s *DeclConstStmt) stmtNode() {}

func (s *DeclConstStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *DeclConstStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())
	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.EQUALS])
	buf.WriteString(" ")
	buf.WriteString(s.Value)

	return buf.String()
}

func (s *DeclConstStmt) Int() int {
	v, err := strconv.Atoi(s.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type ResetStmt struct {
	Kind   lexer.Token // lexer.RESET
	Target []IdentExpr
}

func (s *ResetStmt) stmtNode() {}

func (s *ResetStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ResetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	for i, t := range s.Target {
		buf.WriteString(t.String())

		if len(s.Target)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type ApplyStmt struct {
	Kind   lexer.Token // lexer.X, lexer.CX, ...
	Target []IdentExpr
}

func (s *ApplyStmt) stmtNode() {}

func (s *ApplyStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ApplyStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")

	for i, t := range s.Target {
		buf.WriteString(t.String())

		if len(s.Target)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type MeasureStmt struct {
	Kind   lexer.Token // lexer.MEASURE
	Target []IdentExpr
}

func (s *MeasureStmt) stmtNode() {}

func (s *MeasureStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *MeasureStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	for i, t := range s.Target {
		buf.WriteString(t.String())

		if len(s.Target)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type AssignStmt struct {
	Kind  lexer.Token // lexer.EQUALS
	Left  *IdentExpr
	Right Stmt
}

func (s *AssignStmt) stmtNode() {}

func (s *AssignStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *AssignStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Left.String())
	buf.WriteString(" ")
	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Right.String())

	return buf.String()
}

type PrintStmt struct {
	Kind lexer.Token // lexer.PRINT
}

func (s *PrintStmt) stmtNode() {}

func (s *PrintStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *PrintStmt) String() string {
	return "print"
}
