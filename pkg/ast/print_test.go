package ast_test

import (
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExamplePrint() {
	p := &ast.OpenQASM{
		Version: "3.0",
		Includes: []string{
			"\"stdgates.qasm\"",
		},
		Statements: []ast.Stmt{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.IDENT,
					Value: "q",
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "q",
					},
				},
			},
			&ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.IdentExpr{
					Kind:  lexer.IDENT,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					QArgs: []ast.IdentExpr{
						{
							Kind:  lexer.IDENT,
							Value: "q",
						},
					},
				},
			},
			&ast.ArrowStmt{
				Kind: lexer.ARROW,
				Left: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					QArgs: []ast.IdentExpr{
						{
							Kind:  lexer.IDENT,
							Value: "q",
						},
					},
				},
				Right: &ast.IdentExpr{
					Kind:  lexer.IDENT,
					Value: "c",
				},
			},
		},
	}

	ast.Print(p)

	// Output:
	// *ast.OpenQASM {
	// .  Version: 3.0
	// .  Includes: []string (len = 1) {
	// .  .  0: "stdgates.qasm"
	// .  }
	// .  Statements: []ast.Stmt (len = 5) {
	// .  .  0: *ast.DeclStmt {
	// .  .  .  Kind: qubit
	// .  .  .  Name: *ast.IdentExpr {
	// .  .  .  .  Kind: IDENT
	// .  .  .  .  Value: q
	// .  .  .  }
	// .  .  }
	// .  .  1: *ast.ResetStmt {
	// .  .  .  Kind: reset
	// .  .  .  QArgs: []ast.IdentExpr (len = 1) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  Value: q
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  2: *ast.ApplyStmt {
	// .  .  .  Kind: x
	// .  .  .  QArgs: []ast.IdentExpr (len = 1) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  Value: q
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  3: *ast.AssignStmt {
	// .  .  .  Kind: =
	// .  .  .  Left: *ast.IdentExpr {
	// .  .  .  .  Kind: IDENT
	// .  .  .  .  Value: c
	// .  .  .  }
	// .  .  .  Right: *ast.MeasureStmt {
	// .  .  .  .  Kind: measure
	// .  .  .  .  QArgs: []ast.IdentExpr (len = 1) {
	// .  .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  .  Value: q
	// .  .  .  .  .  }
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  4: *ast.ArrowStmt {
	// .  .  .  Kind: ->
	// .  .  .  Left: *ast.MeasureStmt {
	// .  .  .  .  Kind: measure
	// .  .  .  .  QArgs: []ast.IdentExpr (len = 1) {
	// .  .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  .  Value: q
	// .  .  .  .  .  }
	// .  .  .  .  }
	// .  .  .  }
	// .  .  .  Right: *ast.IdentExpr {
	// .  .  .  .  Kind: IDENT
	// .  .  .  .  Value: c
	// .  .  .  }
	// .  .  }
	// .  }
	// }
}

func ExamplePrint_gate() {
	g := &ast.GateStmt{
		Kind: lexer.GATE,
		Name: "bell",
		QArgs: []ast.IdentExpr{
			{
				Kind:  lexer.IDENT,
				Value: "q0",
			},
			{
				Kind:  lexer.IDENT,
				Value: "q1",
			},
		},
		Statements: []ast.Stmt{
			&ast.ApplyStmt{
				Kind: lexer.H,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "q0",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.CX,
				QArgs: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "q0",
					},
					{
						Kind:  lexer.IDENT,
						Value: "q1",
					},
				},
			},
		},
	}

	ast.Print(g)

	// Output:
	// *ast.GateStmt {
	// .  Kind: gate
	// .  Name: bell
	// .  QArgs: []ast.IdentExpr (len = 2) {
	// .  .  0: ast.IdentExpr {
	// .  .  .  Kind: IDENT
	// .  .  .  Value: q0
	// .  .  }
	// .  .  1: ast.IdentExpr {
	// .  .  .  Kind: IDENT
	// .  .  .  Value: q1
	// .  .  }
	// .  }
	// .  Statements: []ast.Stmt (len = 2) {
	// .  .  0: *ast.ApplyStmt {
	// .  .  .  Kind: h
	// .  .  .  QArgs: []ast.IdentExpr (len = 1) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  Value: q0
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  1: *ast.ApplyStmt {
	// .  .  .  Kind: cx
	// .  .  .  QArgs: []ast.IdentExpr (len = 2) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  Value: q0
	// .  .  .  .  }
	// .  .  .  .  1: ast.IdentExpr {
	// .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  Value: q1
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  }
	// }
}
