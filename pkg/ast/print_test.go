package ast_test

import (
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExamplePrint() {
	p := &ast.OpenQASM{
		Version: "3.0",
		Include: []ast.Expr{
			&ast.IncludeExpr{
				Kind:  lexer.STRING,
				Value: "\"stdgates.qasm\"",
			},
		},
		Statement: []ast.Stmt{
			&ast.DeclStmt{
				Kind: lexer.QUBIT,
				Name: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "q",
				},
			},
			&ast.ResetStmt{
				Kind: lexer.RESET,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.X,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.STRING,
						Value: "q",
					},
				},
			},
			&ast.AssignStmt{
				Kind: lexer.EQUALS,
				Left: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
				Right: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Target: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
				},
			},
			&ast.ArrowStmt{
				Kind: lexer.ARROW,
				Left: &ast.MeasureStmt{
					Kind: lexer.MEASURE,
					Target: []ast.IdentExpr{
						{
							Kind:  lexer.STRING,
							Value: "q",
						},
					},
				},
				Right: &ast.IdentExpr{
					Kind:  lexer.STRING,
					Value: "c",
				},
			},
		},
	}

	ast.Print(p)

	// Output:
	// *ast.OpenQASM {
	// .  Version: 3.0
	// .  Include: []ast.Expr (len = 1) {
	// .  .  0: *ast.IncludeExpr {
	// .  .  .  Kind: STRING
	// .  .  .  Value: "stdgates.qasm"
	// .  .  }
	// .  }
	// .  Statement: []ast.Stmt (len = 5) {
	// .  .  0: *ast.DeclStmt {
	// .  .  .  Kind: qubit
	// .  .  .  Name: *ast.IdentExpr {
	// .  .  .  .  Kind: STRING
	// .  .  .  .  Value: q
	// .  .  .  }
	// .  .  }
	// .  .  1: *ast.ResetStmt {
	// .  .  .  Kind: reset
	// .  .  .  Target: []ast.IdentExpr (len = 1) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: STRING
	// .  .  .  .  .  Value: q
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  2: *ast.ApplyStmt {
	// .  .  .  Kind: x
	// .  .  .  Target: []ast.IdentExpr (len = 1) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: STRING
	// .  .  .  .  .  Value: q
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  3: *ast.AssignStmt {
	// .  .  .  Kind: =
	// .  .  .  Left: *ast.IdentExpr {
	// .  .  .  .  Kind: STRING
	// .  .  .  .  Value: c
	// .  .  .  }
	// .  .  .  Right: *ast.MeasureStmt {
	// .  .  .  .  Kind: measure
	// .  .  .  .  Target: []ast.IdentExpr (len = 1) {
	// .  .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  .  Kind: STRING
	// .  .  .  .  .  .  Value: q
	// .  .  .  .  .  }
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  4: *ast.ArrowStmt {
	// .  .  .  Kind: ->
	// .  .  .  Left: *ast.MeasureStmt {
	// .  .  .  .  Kind: measure
	// .  .  .  .  Target: []ast.IdentExpr (len = 1) {
	// .  .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  .  Kind: STRING
	// .  .  .  .  .  .  Value: q
	// .  .  .  .  .  }
	// .  .  .  .  }
	// .  .  .  }
	// .  .  .  Right: *ast.IdentExpr {
	// .  .  .  .  Kind: STRING
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
		QArg: []ast.IdentExpr{
			{
				Kind:  lexer.STRING,
				Value: "q0",
			},
			{
				Kind:  lexer.STRING,
				Value: "q1",
			},
		},
		Statement: []ast.Stmt{
			&ast.ApplyStmt{
				Kind: lexer.H,
				Target: []ast.IdentExpr{
					{
						Kind:  lexer.IDENT,
						Value: "q0",
					},
				},
			},
			&ast.ApplyStmt{
				Kind: lexer.CX,
				Target: []ast.IdentExpr{
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
	// .  QArg: []ast.IdentExpr (len = 2) {
	// .  .  0: ast.IdentExpr {
	// .  .  .  Kind: STRING
	// .  .  .  Value: q0
	// .  .  }
	// .  .  1: ast.IdentExpr {
	// .  .  .  Kind: STRING
	// .  .  .  Value: q1
	// .  .  }
	// .  }
	// .  Statement: []ast.Stmt (len = 2) {
	// .  .  0: *ast.ApplyStmt {
	// .  .  .  Kind: h
	// .  .  .  Target: []ast.IdentExpr (len = 1) {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: IDENT
	// .  .  .  .  .  Value: q0
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  1: *ast.ApplyStmt {
	// .  .  .  Kind: cx
	// .  .  .  Target: []ast.IdentExpr (len = 2) {
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
