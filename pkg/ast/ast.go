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
