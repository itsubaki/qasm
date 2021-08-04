package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type Program struct {
	Statements []Stmt
}

func (p *Program) String() string {
	var buf bytes.Buffer

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

type Ident struct {
	Kind  lexer.Token // lexer.STRING, lexer.INT, lexer.FLOAT
	Value string
	Index *Index
}

func (i *Ident) stmtNode() {}

func (i *Ident) Token() string {
	return lexer.Tokens[i.Kind]
}

func (i *Ident) String() string {
	if i.Index == nil {
		return i.Value
	}

	return fmt.Sprintf("%s%s", i.Value, i.Index.String())
}

func (i *Ident) IndexValue() int {
	if i.Index == nil {
		return -1
	}

	v, err := strconv.Atoi(i.Index.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type Index struct {
	Kind  lexer.Token // lexer.INT
	Value string
}

func (i *Index) stmtNode() {}

func (i *Index) Token() string {
	return lexer.Tokens[i.Kind]
}

func (i *Index) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[lexer.LBRACKET])
	buf.WriteString(i.Value)
	buf.WriteString(lexer.Tokens[lexer.RBRACKET])

	return buf.String()
}

type LetStmt struct {
	Kind lexer.Token // lexer.QUBIT, lexer.BIT
	Name *Ident
}

func (s *LetStmt) stmtNode() {}

func (s *LetStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *LetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())

	return buf.String()
}

type ResetStmt struct {
	Kind lexer.Token // lexer.RESET
	Name []Ident
}

func (s *ResetStmt) stmtNode() {}

func (s *ResetStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *ResetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
	buf.WriteString(" ")
	for i, e := range s.Name {
		buf.WriteString(e.String())

		if len(s.Name)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type ApplyStmt struct {
	Kind lexer.Token // lexer.X, lexer.CX, ...
	Name *Ident
}

func (s *ApplyStmt) stmtNode() {}

func (s *ApplyStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *ApplyStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())

	return buf.String()
}

type MeasureStmt struct {
	Kind lexer.Token // lexer.MEASURE
	Name *Ident
}

func (s *MeasureStmt) stmtNode() {}

func (s *MeasureStmt) Token() string {
	return lexer.Tokens[s.Kind]
}

func (s *MeasureStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Token())
	buf.WriteString(" ")
	buf.WriteString(s.Name.String())

	return buf.String()
}

type AssignStmt struct {
	Kind  lexer.Token // lexer.EQUALS
	Left  *Ident
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
