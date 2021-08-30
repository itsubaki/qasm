package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type BadExpr struct{}

func (e *BadExpr) exprNode() {}

func (e *BadExpr) Literal() string {
	return ""
}

func (e *BadExpr) String() string {
	return ""
}

type ExprList struct {
	List []Expr
}

func (l *ExprList) Append(x Expr) {
	l.List = append(l.List, x)
}

func (l *ExprList) String() string {
	list := make([]string, 0)
	for _, x := range l.List {
		list = append(list, x.String())
	}

	return strings.Join(list, ", ")
}

type IdentExpr struct {
	Value string
}

func (e *IdentExpr) exprNode() {}

func (e *IdentExpr) Literal() string {
	return lexer.Tokens[lexer.IDENT]
}

func (e *IdentExpr) String() string {
	return e.Value
}

type IndexExpr struct {
	Name  IdentExpr
	Value string
}

func (e *IndexExpr) exprNode() {}

func (e *IndexExpr) Literal() string {
	return lexer.Tokens[lexer.IDENT]
}

func (e *IndexExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(e.Name.Value)
	buf.WriteString(lexer.Tokens[lexer.LBRACKET])
	buf.WriteString(e.Value)
	buf.WriteString(lexer.Tokens[lexer.RBRACKET])

	return buf.String()
}

func (e *IndexExpr) Int() int {
	v, err := strconv.Atoi(e.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type ArrayExpr struct {
	Type IndexExpr
	Name string
}

func (e *ArrayExpr) exprNode() {}

func (e *ArrayExpr) Literal() string {
	return e.Type.Literal()
}

func (e *ArrayExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(e.Type.String())
	buf.WriteString(" ")
	buf.WriteString(e.Name)

	return buf.String()
}

type BasicLit struct {
	Kind  lexer.Token
	Value string
}

func (e *BasicLit) exprNode() {}

func (e *BasicLit) Literal() string {
	return lexer.Tokens[e.Kind]
}

func (e *BasicLit) String() string {
	return e.Value
}

func (e *BasicLit) Float64() float64 {
	v, err := strconv.ParseFloat(e.Value, 64)
	if err != nil {
		panic(err)
	}

	return v
}

type MeasureExpr struct {
	QArgs ExprList
}

func (s *MeasureExpr) exprNode() {}

func (s *MeasureExpr) Literal() string {
	return lexer.Tokens[lexer.MEASURE]
}

func (s *MeasureExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())

	return buf.String()
}

type CallExpr struct {
	Name   string
	Params ParenExpr
	QArgs  ExprList
}

func (e *CallExpr) exprNode() {}

func (e *CallExpr) Literal() string {
	return lexer.Tokens[lexer.GATE]
}

func (e *CallExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(e.Name)
	if len(e.Params.List.List) > 0 {
		buf.WriteString(e.Params.String())
	}

	buf.WriteString(" ")
	buf.WriteString(e.QArgs.String())

	return buf.String()
}

type ParenExpr struct {
	List ExprList
}

func (e *ParenExpr) exprNode() {}

func (e *ParenExpr) Literal() string {
	return lexer.Tokens[lexer.LPAREN]
}

func (e *ParenExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[lexer.LPAREN])
	buf.WriteString(e.List.String())
	buf.WriteString(lexer.Tokens[lexer.RPAREN])

	return buf.String()
}
