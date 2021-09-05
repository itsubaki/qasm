package ast

import (
	"bytes"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type DeclList struct {
	List []Decl
}

func (l *DeclList) Append(d Decl) {
	l.List = append(l.List, d)
}

func (l *DeclList) String() string {
	list := make([]string, 0)
	for _, d := range l.List {
		list = append(list, d.String())
	}

	return strings.Join(list, ", ")
}

type ParenDecl struct {
	List DeclList
}

func (d *ParenDecl) declNode() {}

func (d *ParenDecl) Literal() string {
	return lexer.Tokens[lexer.LPAREN]
}

func (d *ParenDecl) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[lexer.LPAREN])
	buf.WriteString(d.List.String())
	buf.WriteString(lexer.Tokens[lexer.RPAREN])

	return buf.String()
}

func (d *ParenDecl) Append(decl Decl) {
	d.List.Append(decl)
}

type BadDecl struct{}

func (d *BadDecl) declNode() {}

func (d *BadDecl) Literal() string {
	return ""
}

func (d *BadDecl) String() string {
	return ""
}

type GenDecl struct {
	Kind lexer.Token // lexer.QUBIT, lexer.BIT
	Type Expr
	Name IdentExpr
}

func (d *GenDecl) declNode() {}

func (d *GenDecl) Literal() string {
	return lexer.Tokens[d.Kind]
}

func (d *GenDecl) String() string {
	var buf bytes.Buffer

	buf.WriteString(d.Type.String())
	buf.WriteString(" ")
	buf.WriteString(d.Name.Value)

	return buf.String()
}

func (d *GenDecl) Size() int {
	switch x := d.Type.(type) {
	case *IndexExpr:
		return x.Int()
	}

	return 1
}

type GenConst struct {
	Name  IdentExpr
	Value Expr
}

func (d *GenConst) declNode() {}

func (d *GenConst) Literal() string {
	return lexer.Tokens[lexer.CONST]
}

func (d *GenConst) String() string {
	var buf bytes.Buffer

	buf.WriteString(d.Literal())
	buf.WriteString(" ")
	buf.WriteString(d.Name.String())
	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.EQUALS])
	buf.WriteString(" ")
	buf.WriteString(d.Value.String())

	return buf.String()
}

type GateDecl struct {
	Name   string
	Params ParenExpr
	QArgs  ExprList
	Body   BlockStmt
}

func (d *GateDecl) declNode() {}

func (d *GateDecl) Literal() string {
	return lexer.Tokens[lexer.GATE]
}

func (d *GateDecl) String() string {
	var buf bytes.Buffer

	buf.WriteString(d.Literal())
	buf.WriteString(" ")

	buf.WriteString(d.Name)
	if len(d.Params.List.List) > 0 {
		buf.WriteString(d.Params.String())
	}

	buf.WriteString(" ")
	buf.WriteString(d.QArgs.String())
	buf.WriteString(" ")
	buf.WriteString(d.Body.String())

	return buf.String()
}

type FuncDecl struct {
	Name   string
	Params ParenDecl
	QArgs  DeclList
	Body   BlockStmt
	Result Expr
}

func (d *FuncDecl) declNode() {}

func (d *FuncDecl) Literal() string {
	return lexer.Tokens[lexer.DEF]
}

func (d *FuncDecl) String() string {
	var buf bytes.Buffer

	buf.WriteString(d.Literal())
	buf.WriteString(" ")

	buf.WriteString(d.Name)
	if len(d.Params.List.List) > 0 {
		buf.WriteString(d.Params.String())
	}

	buf.WriteString(" ")
	buf.WriteString(d.QArgs.String())

	if d.Result != nil {
		buf.WriteString(" ")
		buf.WriteString(lexer.Tokens[lexer.ARROW])

		buf.WriteString(" ")
		buf.WriteString(d.Result.String())
	}

	buf.WriteString(" ")
	buf.WriteString(d.Body.String())

	return buf.String()
}

type VersionDecl struct {
	Value Expr
}

func (d *VersionDecl) declNode() {}

func (d *VersionDecl) Literal() string {
	return lexer.Tokens[lexer.OPENQASM]
}

func (d *VersionDecl) String() string {
	var buf bytes.Buffer

	buf.WriteString(d.Literal())
	buf.WriteString(" ")
	buf.WriteString(d.Value.String())

	return buf.String()
}
