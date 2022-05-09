package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type BadExpr struct{}

func (x *BadExpr) exprNode() {}

func (x *BadExpr) Literal() string {
	return ""
}

func (x *BadExpr) String() string {
	return ""
}

type ExprList struct {
	List []Expr
}

func (x *ExprList) exprNode() {}

func (x *ExprList) Literal() string {
	return ""
}

func (x *ExprList) Append(e Expr) {
	x.List = append(x.List, e)
}

func (x *ExprList) String() string {
	var list []string
	for _, e := range x.List {
		list = append(list, e.String())
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

type IdentExpr struct {
	Name string
}

func (x *IdentExpr) exprNode() {}

func (x *IdentExpr) Literal() string {
	return lexer.Tokens[lexer.IDENT]
}

func (x *IdentExpr) String() string {
	return x.Name
}

type IndexExpr struct {
	Name  string
	Value string
}

func (x *IndexExpr) exprNode() {}

func (x *IndexExpr) Literal() string {
	return lexer.Tokens[lexer.IDENT]
}

func (x *IndexExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Name)
	buf.WriteString(lexer.Tokens[lexer.LBRACKET])
	buf.WriteString(x.Value)
	buf.WriteString(lexer.Tokens[lexer.RBRACKET])

	return buf.String()
}

func (x *IndexExpr) Int() int {
	return Must(strconv.Atoi(x.Value))
}

type BasicLit struct {
	Kind  lexer.Token // lexer.INT, lexer.FLOAT, lexer.STRING
	Value string
}

func (x *BasicLit) exprNode() {}

func (x *BasicLit) Literal() string {
	return lexer.Tokens[x.Kind]
}

func (x *BasicLit) String() string {
	return x.Value
}

func (x *BasicLit) Int64() int64 {
	return Must(strconv.ParseInt(x.Value, 10, 64))
}

func (x *BasicLit) Float64() float64 {
	return Must(strconv.ParseFloat(x.Value, 64))
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

type Modifier struct {
	Kind  lexer.Token // lexer.CTRL, lexer.NEGCTRL, lexer.INV, lexer.POW
	Index ParenExpr
}

func (x *Modifier) exprNode() {}

func (x *Modifier) Literal() string {
	return lexer.Tokens[x.Kind]
}

func (x *Modifier) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Literal())
	if len(x.Index.List.List) > 0 {
		buf.WriteString(x.Index.String())
	}

	return buf.String()
}

func ModPow(mod []Modifier) []Modifier {
	var out []Modifier
	for i, m := range mod {
		if m.Kind == lexer.POW {
			out = append(out, mod[i])
		}
	}

	return out
}

func ModInv(mod []Modifier) []Modifier {
	var out []Modifier
	for i, m := range mod {
		if m.Kind == lexer.INV {
			out = append(out, mod[i])
		}
	}

	return out
}

func ModCtrl(mod []Modifier) []Modifier {
	var out []Modifier
	for i, m := range mod {
		if m.Kind == lexer.CTRL || m.Kind == lexer.NEGCTRL {
			out = append(out, mod[i])
		}
	}

	return out
}

type CallExpr struct {
	Name     string
	Modifier []Modifier // lexer.CTRL, lexer.NEGCTRL, lexer.INV, lexer.POW
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

type InfixExpr struct {
	Kind  lexer.Token
	Left  Expr
	Right Expr
}

func (x *InfixExpr) exprNode() {}

func (x *InfixExpr) Literal() string {
	return lexer.Tokens[x.Kind]
}

func (x *InfixExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Left.String())
	buf.WriteString(" ")
	buf.WriteString(x.Literal())
	buf.WriteString(" ")
	buf.WriteString(x.Right.String())

	return buf.String()
}

type UnaryExpr struct {
	Kind  lexer.Token
	Value Expr
}

func (x *UnaryExpr) exprNode() {}

func (x *UnaryExpr) Literal() string {
	return lexer.Tokens[x.Kind]
}

func (x *UnaryExpr) String() string {
	var buf bytes.Buffer

	buf.WriteString(x.Literal())
	buf.WriteString(x.Value.String())

	return buf.String()
}
