package ast

import (
	"bytes"
	"fmt"
)

type OpenQASM struct {
	Version string
	Incls   []Stmt
	Stmts   []Stmt
}

func (p *OpenQASM) String() string {
	var buf bytes.Buffer

	version := fmt.Sprintf("OPENQASM %v;\n", p.Version)
	buf.WriteString(version)

	for _, s := range p.Incls {
		str := fmt.Sprintf("%s\n", s.String())
		buf.WriteString(str)
	}

	for _, s := range p.Stmts {
		str := fmt.Sprintf("%s\n", s.String())
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

type Decl interface {
	Node
	declNode()
}

func Ident(x interface{}) string {
	switch x := x.(type) {
	case *IdentExpr:
		return x.Value
	case *IndexExpr:
		return x.Name.Value
	case *ArrayExpr:
		return x.Name
	case *GenDecl:
		return x.Name.Value
	case *GenConst:
		return x.Name.Value
	case *GateDecl:
		return x.Name
	case *FuncDecl:
		return x.Name
	case *BasicExpr:
		return x.Value
	}

	panic(fmt.Errorf("invalid type=%#v", x))
}
