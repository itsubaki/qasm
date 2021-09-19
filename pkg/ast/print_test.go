package ast_test

import (
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func Example() {
	p := &ast.OpenQASM{
		Version: &ast.DeclStmt{
			Decl: &ast.VersionDecl{
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3.0",
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.InclStmt{
				Path: ast.BasicLit{
					Kind:  lexer.STRING,
					Value: "\"testdata/stdgates.qasm\"",
				},
			},
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Kind: lexer.QUBIT,
					Type: &ast.IdentExpr{
						Value: lexer.Tokens[lexer.QUBIT],
					},
					Name: ast.IdentExpr{
						Value: "q",
					},
				},
			},
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Value: "q",
						},
					},
				},
			},
		},
	}

	ast.Println(p)

	// Output:
	// *ast.OpenQASM {
	// .  Version: *ast.DeclStmt {
	// .  .  Decl: *ast.VersionDecl {
	// .  .  .  Value: *ast.BasicLit {
	// .  .  .  .  Kind: FLOAT
	// .  .  .  .  Value: 3.0
	// .  .  .  }
	// .  .  }
	// .  }
	// .  Stmts: []ast.Stmt (len = 3) {
	// .  .  0: *ast.InclStmt {
	// .  .  .  Path: ast.BasicLit {
	// .  .  .  .  Kind: STRING
	// .  .  .  .  Value: "testdata/stdgates.qasm"
	// .  .  .  }
	// .  .  }
	// .  .  1: *ast.DeclStmt {
	// .  .  .  Decl: *ast.GenDecl {
	// .  .  .  .  Kind: qubit
	// .  .  .  .  Type: *ast.IdentExpr {
	// .  .  .  .  .  Value: qubit
	// .  .  .  .  }
	// .  .  .  .  Name: ast.IdentExpr {
	// .  .  .  .  .  Value: q
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  2: *ast.ResetStmt {
	// .  .  .  QArgs: ast.ExprList {
	// .  .  .  .  List: []ast.Expr (len = 1) {
	// .  .  .  .  .  0: *ast.IdentExpr {
	// .  .  .  .  .  .  Value: q
	// .  .  .  .  .  }
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  }
	// }
}
