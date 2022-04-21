package ast

import (
	"bytes"
	"fmt"
)

type OpenQASM struct {
	Version Stmt
	Stmts   []Stmt
}

func (p *OpenQASM) String() string {
	var buf bytes.Buffer

	v := fmt.Sprintf("%v\n", p.Version.String())
	buf.WriteString(v)

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
		return x.Name
	case *IndexExpr:
		return x.Name
	case *GenDecl:
		return x.Name
	case *GenConst:
		return x.Name
	case *GateDecl:
		return x.Name
	case *FuncDecl:
		return x.Name
	case *BasicLit:
		return x.Value
	}

	panic(fmt.Sprintf("invalid type=%T", x))
}
