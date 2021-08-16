package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type OpenQASM struct {
	Version   string
	Include   []Expr
	Gate      []Stmt
	Statement []Stmt
}

func (p *OpenQASM) String() string {
	var buf bytes.Buffer

	version := fmt.Sprintf("OPENQASM %v;\n", p.Version)
	buf.WriteString(version)

	for _, i := range p.Include {
		str := fmt.Sprintf("include %s;\n", i)
		buf.WriteString(str)
	}

	for _, g := range p.Gate {
		str := fmt.Sprintf("%s\n", g.String())
		buf.WriteString(str)
	}

	for _, s := range p.Statement {
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

type IncludeExpr struct {
	Kind  lexer.Token // lexer.STRING
	Value string
}

func (i *IncludeExpr) exprNode() {}

func (i *IncludeExpr) Literal() string {
	return lexer.Tokens[i.Kind]
}

func (i *IncludeExpr) String() string {
	return i.Value
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
	Kind  lexer.Token // lexer.INT
	Value string
}

func (i *IndexExpr) exprNode() {}

func (i *IndexExpr) Literal() string {
	return lexer.Tokens[i.Kind]
}

func (i *IndexExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[lexer.LBRACKET])
	buf.WriteString(i.Value)
	buf.WriteString(lexer.Tokens[lexer.RBRACKET])

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
	Kind  lexer.Token // lexer.QUBIT, lexer.BIT, lexer.CONST
	Name  *IdentExpr
	Value string
}

func (s *DeclStmt) stmtNode() {}

func (s *DeclStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *DeclStmt) String() string {
	var buf bytes.Buffer

	if s.Kind == lexer.CONST {
		buf.WriteString(s.Literal())
		buf.WriteString(" ")
		buf.WriteString(s.Name.String())
		buf.WriteString(" ")
		buf.WriteString(lexer.Tokens[lexer.EQUALS])
		buf.WriteString(" ")
		buf.WriteString(s.Value)

		return buf.String()
	}

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())

	return buf.String()
}

func (s *DeclStmt) Int() int {
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

type ArrowStmt struct {
	Kind  lexer.Token // lexer.ARROW
	Left  Stmt
	Right *IdentExpr
}

func (s *ArrowStmt) stmtNode() {}

func (s *ArrowStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ArrowStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Left.String())
	buf.WriteString(" ")
	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Right.String())

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
	Kind   lexer.Token // lexer.PRINT
	Target []IdentExpr
}

func (s *PrintStmt) stmtNode() {}

func (s *PrintStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *PrintStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if s.Target == nil || len(s.Target) == 0 {
		return buf.String()
	}

	buf.WriteString(" ")
	for i, t := range s.Target {
		buf.WriteString(t.String())

		if len(s.Target)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type GateStmt struct {
	Kind      lexer.Token // lexer.GATE
	Name      string
	QArg      []IdentExpr
	Statement []Stmt
}

func (s *GateStmt) stmtNode() {}

func (s *GateStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *GateStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Name)
	buf.WriteString(" ")
	for i, a := range s.QArg {
		buf.WriteString(a.String())

		if len(s.QArg)-1 != i {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.LBRACE])
	buf.WriteString(" ")
	for _, stmt := range s.Statement {
		msg := fmt.Sprintf("%s; ", stmt.String())
		buf.WriteString(msg)
	}
	buf.WriteString(lexer.Tokens[lexer.RBRACE])

	return buf.String()
}
