package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

func Ident(x Expr) string {
	switch x := x.(type) {
	case *IdentExpr:
		return x.Value
	case *IndexExpr:
		return x.Name.Value
	case *ArrayExpr:
		return x.Name
	}

	return x.String()
}

type ExprList struct {
	List []Expr
}

func (l *ExprList) Append(x Expr) {
	l.List = append(l.List, x)
}

func (l *ExprList) String() string {
	list := make([]string, 0)
	for _, t := range l.List {
		list = append(list, t.String())
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
	Name  *IdentExpr
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
	Type *IndexExpr
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

type ResetExpr struct {
	QArgs ExprList
}

func (s *ResetExpr) exprNode() {}

func (s *ResetExpr) Literal() string {
	return lexer.Tokens[lexer.RESET]
}

func (s *ResetExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())

	return buf.String()
}

type PrintExpr struct {
	QArgs ExprList
}

func (s *PrintExpr) exprNode() {}

func (s *PrintExpr) Literal() string {
	return lexer.Tokens[lexer.PRINT]
}

func (s *PrintExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if len(s.QArgs.List) == 0 {
		return buf.String()
	}

	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())

	return buf.String()
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

type ApplyExpr struct {
	Kind   lexer.Token // lexer.X, lexer.CX, ...
	Params ExprList
	QArgs  ExprList
}

func (s *ApplyExpr) exprNode() {}

func (s *ApplyExpr) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ApplyExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if len(s.Params.List) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		buf.WriteString(s.Params.String())
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())

	return buf.String()
}

type CallExpr struct {
	Name   string
	Params ExprList
	QArgs  ExprList
}

func (e *CallExpr) exprNode() {}

func (e *CallExpr) Literal() string {
	return lexer.Tokens[lexer.GATE]
}

func (e *CallExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(e.Name)
	if len(e.Params.List) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		buf.WriteString(e.Params.String())
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	buf.WriteString(e.QArgs.String())

	return buf.String()
}
