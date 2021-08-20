package ast

import (
	"bytes"
	"strconv"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type GenDecl struct {
	Kind lexer.Token // lexer.QUBIT, lexer.BIT
	Type Expr
	Name *IdentExpr
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
	Name  *IdentExpr
	Value string
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
	buf.WriteString(d.Value)

	return buf.String()
}

func (s *GenConst) Int() int {
	v, err := strconv.Atoi(s.Value)
	if err != nil {
		panic(err)
	}

	return v
}

type FuncDecl struct {
	Kind   lexer.Token // lexer.GATE
	Name   string
	Params ExprList
	QArgs  ExprList
	Body   *BlockStmt
	Output *IndexExpr
}

func (d *FuncDecl) declNode() {}

func (d *FuncDecl) Literal() string {
	return lexer.Tokens[d.Kind]
}

func (d *FuncDecl) String() string {
	var buf bytes.Buffer

	buf.WriteString(d.Literal())
	buf.WriteString(" ")
	buf.WriteString(d.Name)
	if len(d.Params.List) > 0 {
		buf.WriteString(lexer.Tokens[lexer.LPAREN])
		buf.WriteString(d.Params.String())
		buf.WriteString(lexer.Tokens[lexer.RPAREN])
	}

	buf.WriteString(" ")
	buf.WriteString(d.QArgs.String())
	buf.WriteString(" ")
	if d.Output != nil {
		buf.WriteString("->")
		buf.WriteString(" ")
		buf.WriteString(d.Output.String())
		buf.WriteString(" ")
	}
	buf.WriteString(d.Body.String())

	return buf.String()
}
