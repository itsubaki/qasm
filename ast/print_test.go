package ast_test

import (
	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/lexer"
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
						Name: lexer.Tokens[lexer.QUBIT],
					},
					Name: "q",
				},
			},
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Name: "q",
						},
					},
				},
			},
			nil,
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
	// .  Stmts: []ast.Stmt (len = 4) {
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
	// .  .  .  .  .  Name: qubit
	// .  .  .  .  }
	// .  .  .  .  Name: q
	// .  .  .  }
	// .  .  }
	// .  .  2: *ast.ResetStmt {
	// .  .  .  QArgs: ast.ExprList {
	// .  .  .  .  List: []ast.Expr (len = 1) {
	// .  .  .  .  .  0: *ast.IdentExpr {
	// .  .  .  .  .  .  Name: q
	// .  .  .  .  .  }
	// .  .  .  .  }
	// .  .  .  }
	// .  .  }
	// .  .  3: nil
	// .  }
	// }
}
