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
	Token() string
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

func (i *IdentExpr) Token() string {
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

func (i *IndexExpr) Token() string {
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

type LetStmt struct {
	Kind  lexer.Token // lexer.QUBIT, lexer.BIT
	Name  *IdentExpr
	Index *IndexExpr
}

func (s *LetStmt) stmtNode() {}

func (s *LetStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *LetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
	if s.Index != nil {
		buf.WriteString(s.Index.String())
	}
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())

	return buf.String()
}

type ResetStmt struct {
	Kind   lexer.Token // lexer.RESET
	Target []IdentExpr
}

func (s *ResetStmt) stmtNode() {}

func (s *ResetStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *ResetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
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

func (s *ApplyStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *ApplyStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
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

func (s *MeasureStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *MeasureStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
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

func (s *AssignStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *AssignStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Left.String())
	buf.WriteString(" ")
	buf.WriteString(s.Token())
	buf.WriteString(" ")
	buf.WriteString(s.Right.String())

	return buf.String()
}

type PrintStmt struct {
	Kind lexer.Token // lexer.PRINT
}

func (s *PrintStmt) stmtNode() {}

func (s *PrintStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *PrintStmt) String() string {
	return "print"
}
