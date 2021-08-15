package ast_test

import (
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func ExamplePrint() {
	p := &ast.OpenQASM{
		Version: "3.0",
		Includes: []ast.Expr{
			&ast.IncludeExpr{
				Kind:  lexer.STRING,
				Value: "\"stdgates.qasm\"",
			},
		},
		Statements: []ast.Stmt{
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
	// .  Includes: []ast.Expr {
	// .  .  0: *ast.IncludeExpr {
	// .  .  .  Kind: STRING
	// .  .  .  Value: "stdgates.qasm"
	// .  .  }
	// .  }
	// .  Statements: []ast.Stmt {
	// .  .  0: *ast.DeclStmt {
	// .  .  .  Kind: qubit
	// .  .  .  Name: *ast.IdentExpr {
	// .  .  .  .  Kind: STRING
	// .  .  .  .  Value: q
	// .  .  .  }
	// .  .  }
	// .  .  1: *ast.ResetStmt {
	// .  .  .  Kind: reset
	// .  .  .  Target: []ast.IdentExpr {
	// .  .  .  .  0: ast.IdentExpr {
	// .  .  .  .  .  Kind: STRING
	// .  .  .  .  .  Value: q
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  2: *ast.ApplyStmt {
	// .  .  .  Kind: x
	// .  .  .  Target: []ast.IdentExpr {
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
	// .  .  .  .  Target: []ast.IdentExpr {
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
	// .  .  .  .  Target: []ast.IdentExpr {
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
