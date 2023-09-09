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

func Ident(x interface{}) (string, error) {
	switch x := x.(type) {
	case *IdentExpr:
		return x.Name, nil // q
	case *IndexExpr:
		return x.Name, nil // q[2]
	case *GenDecl:
		return x.Name, nil // qubit q, qubit[2] q
	case *GenConst:
		return x.Name, nil // const N = 15
	case *GateDecl:
		return x.Name, nil // gate X {}
	case *SubroutineDecl:
		return x.Name, nil // def shor(){}
	case *BasicLit:
		return x.Value, nil
	}

	return "", fmt.Errorf("invalid type=%T", x)
}
