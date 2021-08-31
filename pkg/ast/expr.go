package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

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

type ParenExpr struct {
	List ExprList
}

func (x *ParenExpr) exprNode() {}

func (x *ParenExpr) Literal() string {
	return lexer.Tokens[lexer.LPAREN]
}

func (x *ParenExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[lexer.LPAREN])
	buf.WriteString(x.List.String())
	buf.WriteString(lexer.Tokens[lexer.RPAREN])

	return buf.String()
}

type BadExpr struct{}

func (x *BadExpr) exprNode() {}

func (x *BadExpr) Literal() string {
	return ""
}

func (x *BadExpr) String() string {
	return ""
}

type IdentExpr struct {
	Value string
}

func (x *IdentExpr) exprNode() {}

func (x *IdentExpr) Literal() string {
	return lexer.Tokens[lexer.IDENT]
}

func (x *IdentExpr) String() string {
	return x.Value
}

type IndexExpr struct {
	Name  IdentExpr
	Value string
}

func (x *IndexExpr) exprNode() {}

func (x *IndexExpr) Literal() string {
	return lexer.Tokens[lexer.IDENT]
}

func (x *IndexExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Name.Value)
	buf.WriteString(lexer.Tokens[lexer.LBRACKET])
	buf.WriteString(x.Value)
	buf.WriteString(lexer.Tokens[lexer.RBRACKET])

	return buf.String()
}

func (x *IndexExpr) Int() int {
	v, err := strconv.Atoi(x.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type ArrayExpr struct {
	Type IndexExpr
	Name string
}

func (x *ArrayExpr) exprNode() {}

func (x *ArrayExpr) Literal() string {
	return x.Type.Literal()
}

func (x *ArrayExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Type.String())
	buf.WriteString(" ")
	buf.WriteString(x.Name)

	return buf.String()
}

type BasicLit struct {
	Kind  lexer.Token
	Value string
}

func (x *BasicLit) exprNode() {}

func (x *BasicLit) Literal() string {
	return lexer.Tokens[x.Kind]
}

func (x *BasicLit) String() string {
	return x.Value
}

func (x *BasicLit) Float64() float64 {
	v, err := strconv.ParseFloat(x.Value, 64)
	if err != nil {
		panic(err)
	}

	return v
}

type MeasureExpr struct {
	QArgs ExprList
}

func (x *MeasureExpr) exprNode() {}

func (x *MeasureExpr) Literal() string {
	return lexer.Tokens[lexer.MEASURE]
}

func (x *MeasureExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Literal())
	buf.WriteString(" ")
	buf.WriteString(x.QArgs.String())

	return buf.String()
}

type Modifiler struct {
	Kind  lexer.Token // lexer.CTRL, lexer.NEGCTRL, lexer.INV
	Index ParenExpr
}

func (x *Modifiler) exprNode() {}

func (x *Modifiler) Literal() string {
	return lexer.Tokens[x.Kind]
}

func (x *Modifiler) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Literal())
	if len(x.Index.List.List) > 0 {
		buf.WriteString(x.Index.String())
	}

	return buf.String()
}

type CallExpr struct {
	Name     string
	Modifier []Modifiler // lexer.CTRL, lexer.NEGCTRL, lexer.INV
	Params   ParenExpr
	QArgs    ExprList
}

func (x *CallExpr) exprNode() {}

func (x *CallExpr) Literal() string {
	return lexer.Tokens[lexer.GATE]
}

func (x *CallExpr) String() string {
	var buf bytes.Buffer

	for _, m := range x.Modifier {
		buf.WriteString(m.String())
		buf.WriteString(" ")
		buf.WriteString(lexer.Tokens[lexer.AT])
		buf.WriteString(" ")
	}

	buf.WriteString(x.Name)
	if len(x.Params.List.List) > 0 {
		buf.WriteString(x.Params.String())
	}

	buf.WriteString(" ")
	buf.WriteString(x.QArgs.String())

	return buf.String()
}
