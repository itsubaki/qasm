package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type OpenQASM struct {
	Version    string
	Includes   []string
	Gates      []Stmt
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

	for _, g := range p.Gates {
		str := fmt.Sprintf("%s\n", g.String())
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
	Kind  lexer.Token // lexer.IDENT
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

type GateStmt struct {
	Kind       lexer.Token // lexer.GATE
	Name       string
	Params     []IdentExpr
	QArgs      []IdentExpr
	Statements []Stmt
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

	if len(s.Params) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		for i, p := range s.Params {
			buf.WriteString(p.String())

			if len(s.Params)-1 != i {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	for i, a := range s.QArgs {
		buf.WriteString(a.String())

		if len(s.QArgs)-1 != i {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.LBRACE])
	buf.WriteString(" ")
	for _, stmt := range s.Statements {
		msg := fmt.Sprintf("%s; ", stmt.String())
		buf.WriteString(msg)
	}
	buf.WriteString(lexer.Tokens[lexer.RBRACE])

	return buf.String()
}

type DefStmt struct {
	Kind       lexer.Token // lexer.DEF
	Name       string
	Params     []IdentExpr
	QArgs      []IdentExpr
	Output     *IdentExpr
	Statements []Stmt
}

func (s *DefStmt) stmtNode() {}

func (s *DefStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *DefStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Name)

	if len(s.Params) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		for i, p := range s.Params {
			buf.WriteString(p.String())

			if len(s.Params)-1 != i {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	for i, a := range s.QArgs {
		buf.WriteString(lexer.Tokens[lexer.QUBIT])
		buf.WriteString(lexer.Tokens[lexer.LBRACKET])
		buf.WriteString(a.Index.Value)
		buf.WriteString(lexer.Tokens[lexer.RBRACKET])
		buf.WriteString(" ")
		buf.WriteString(a.Value)

		if len(s.QArgs)-1 != i {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.ARROW])
	buf.WriteString(" ")
	buf.WriteString(s.Output.Literal())
	buf.WriteString(lexer.Tokens[lexer.LBRACKET])
	buf.WriteString(s.Output.Index.Value)
	buf.WriteString(lexer.Tokens[lexer.RBRACKET])
	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.LBRACE])
	buf.WriteString(" ")
	for _, stmt := range s.Statements {
		msg := fmt.Sprintf("%s; ", stmt.String())
		buf.WriteString(msg)
	}
	buf.WriteString(lexer.Tokens[lexer.RBRACE])

	return buf.String()
}

type CallStmt struct {
	Kind   lexer.Token // lexer.IDENT
	Name   string
	Params []IdentExpr
	QArgs  []IdentExpr
}

func (s *CallStmt) stmtNode() {}

func (s *CallStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *CallStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Name)
	if len(s.Params) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		for i, p := range s.Params {
			buf.WriteString(p.String())

			if len(s.Params)-1 != i {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	for i, a := range s.QArgs {
		buf.WriteString(a.String())

		if len(s.QArgs)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type ReturnStmt struct {
	Kind  lexer.Token // lexer.RETURN
	Value Stmt
}

func (s *ReturnStmt) stmtNode() {}

func (s *ReturnStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ReturnStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Value.String())

	return buf.String()
}

type ExprStmt struct {
	X Expr
}

type ConstStmt struct {
	Kind  lexer.Token //  lexer.CONST
	Name  *IdentExpr
	Value string
}

func (s *ConstStmt) stmtNode() {}

func (s *ConstStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ConstStmt) String() string {
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

func (s *ConstStmt) Int() int {
	v, err := strconv.Atoi(s.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type DeclStmt struct {
	Kind lexer.Token // lexer.QUBIT, lexer.BIT
	Name *IdentExpr
}

func (s *DeclStmt) stmtNode() {}

func (s *DeclStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *DeclStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if s.Name.Index != nil {
		buf.WriteString(lexer.Tokens[lexer.LBRACKET])
		buf.WriteString(s.Name.Index.Value)
		buf.WriteString(lexer.Tokens[lexer.RBRACKET])
	}
	buf.WriteString(" ")
	buf.WriteString(s.Name.Value)

	return buf.String()
}

type ResetStmt struct {
	Kind  lexer.Token // lexer.RESET
	QArgs []IdentExpr
}

func (s *ResetStmt) stmtNode() {}

func (s *ResetStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ResetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	for i, t := range s.QArgs {
		buf.WriteString(t.String())

		if len(s.QArgs)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type ApplyStmt struct {
	Kind   lexer.Token // lexer.X, lexer.CX, ...
	Params []IdentExpr
	QArgs  []IdentExpr
}

func (s *ApplyStmt) stmtNode() {}

func (s *ApplyStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ApplyStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if len(s.Params) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		for i, p := range s.Params {
			buf.WriteString(p.String())

			if len(s.Params)-1 != i {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	for i, t := range s.QArgs {
		buf.WriteString(t.String())

		if len(s.QArgs)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

type MeasureStmt struct {
	Kind  lexer.Token // lexer.MEASURE
	QArgs []IdentExpr
}

func (s *MeasureStmt) stmtNode() {}

func (s *MeasureStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *MeasureStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	for i, t := range s.QArgs {
		buf.WriteString(t.String())

		if len(s.QArgs)-1 != i {
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
	Kind  lexer.Token // lexer.PRINT
	QArgs []IdentExpr
}

func (s *PrintStmt) stmtNode() {}

func (s *PrintStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *PrintStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if s.QArgs == nil || len(s.QArgs) == 0 {
		return buf.String()
	}

	buf.WriteString(" ")
	for i, t := range s.QArgs {
		buf.WriteString(t.String())

		if len(s.QArgs)-1 != i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}
