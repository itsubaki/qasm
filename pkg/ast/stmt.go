package ast

import (
	"bytes"
	"fmt"

	"github.com/itsubaki/qasm/pkg/lexer"
)

type BadStmt struct{}

func (s *BadStmt) stmtNode() {}

func (s *BadStmt) Literal() string {
	return ""
}

func (s *BadStmt) String() string {
	return ""
}

type ExprStmt struct {
	X Expr
}

func (s *ExprStmt) stmtNode() {}

func (s *ExprStmt) Literal() string {
	return s.X.Literal()
}

func (s *ExprStmt) String() string {
	return fmt.Sprintf("%s;", s.X)
}

type DeclStmt struct {
	Decl Decl
}

func (s *DeclStmt) stmtNode() {}

func (s *DeclStmt) Literal() string {
	return s.Decl.Literal()
}

func (s *DeclStmt) String() string {
	switch d := s.Decl.(type) {
	case *GateDecl, *FuncDecl:
		return d.String()
	}

	return fmt.Sprintf("%s;", s.Decl)
}

type InclStmt struct {
	Path BasicExpr
}

func (s *InclStmt) stmtNode() {}

func (s *InclStmt) Literal() string {
	return lexer.Tokens[lexer.INCLUDE]
}

func (s *InclStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Path.String())
	buf.WriteString(";")

	return buf.String()
}

type ReturnStmt struct {
	Result Expr
}

func (s *ReturnStmt) stmtNode() {}

func (s *ReturnStmt) Literal() string {
	return lexer.Tokens[lexer.RETURN]
}

func (s *ReturnStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if s.Result != nil {
		buf.WriteString(" ")
		buf.WriteString(s.Result.String())
	}
	buf.WriteString(";")

	return buf.String()
}

type ArrowStmt struct {
	Left  Expr
	Right Expr
}

func (s *ArrowStmt) stmtNode() {}

func (s *ArrowStmt) Literal() string {
	return lexer.Tokens[lexer.ARROW]
}

func (s *ArrowStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Left.String())
	buf.WriteString(" ")
	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Right.String())
	buf.WriteString(";")

	return buf.String()
}

type AssignStmt struct {
	Left  Expr
	Right Expr
}

func (s *AssignStmt) stmtNode() {}

func (s *AssignStmt) Literal() string {
	return lexer.Tokens[lexer.EQUALS]
}

func (s *AssignStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Left.String())
	buf.WriteString(" ")
	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.Right.String())
	buf.WriteString(";")

	return buf.String()
}

type ResetStmt struct {
	QArgs ExprList
}

func (s *ResetStmt) stmtNode() {}

func (s *ResetStmt) Literal() string {
	return lexer.Tokens[lexer.RESET]
}

func (s *ResetStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())
	buf.WriteString(";")

	return buf.String()
}

type PrintStmt struct {
	QArgs ExprList
}

func (s *PrintStmt) stmtNode() {}

func (s *PrintStmt) Literal() string {
	return lexer.Tokens[lexer.PRINT]
}

func (s *PrintStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if len(s.QArgs.List) == 0 {
		buf.WriteString(";")
		return buf.String()
	}

	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())
	buf.WriteString(";")

	return buf.String()
}

type ApplyStmt struct {
	Kind   lexer.Token // lexer.X, lexer.CX, ...
	Params ParenExpr
	QArgs  ExprList
}

func (s *ApplyStmt) stmtNode() {}

func (s *ApplyStmt) Literal() string {
	return lexer.Tokens[s.Kind]
}

func (s *ApplyStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(s.Literal())
	if len(s.Params.List.List) > 0 {
		buf.WriteString(s.Params.String())
	}

	buf.WriteString(" ")
	buf.WriteString(s.QArgs.String())
	buf.WriteString(";")

	return buf.String()
}

type BlockStmt struct {
	List []Stmt
}

func (s *BlockStmt) stmtNode() {}

func (s *BlockStmt) Literal() string {
	return lexer.Tokens[lexer.LBRACE]
}

func (s *BlockStmt) String() string {
	var buf bytes.Buffer

	buf.WriteString(lexer.Tokens[lexer.LBRACE])
	for _, s := range s.List {
		buf.WriteString(" ")
		buf.WriteString(s.String())
	}
	buf.WriteString(" ")
	buf.WriteString(lexer.Tokens[lexer.RBRACE])

	return buf.String()
}

func (s *BlockStmt) Append(stmt Stmt) {
	s.List = append(s.List, stmt)
}

type IfStmt struct{}

type BranchStmt struct{}

type ForStmt struct{}
