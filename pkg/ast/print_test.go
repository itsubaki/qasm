package ast_test

// import (
// 	"github.com/itsubaki/qasm/pkg/ast"
// 	"github.com/itsubaki/qasm/pkg/lexer"
// )

// func ExamplePrint() {
// 	p := &ast.OpenQASM{
// 		Version: "3.0",
// 		Incls: []string{
// 			"\"stdgates.qasm\"",
// 		},
// 		Stmts: []ast.Stmt{
// 			&ast.DeclStmt{
// 				Decl: &ast.GenDecl{
// 					Kind: lexer.QUBIT,
// 					Name: &ast.IdentExpr{
// 						Value: "q",
// 					},
// 				},
// 			},
// 			&ast.ExprStmt{
// 				X: &ast.ResetExpr{
// 					QArgs: ast.IdentList{
// 						List: []ast.IdentExpr{
// 							{
// 								Value: "q",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			&ast.ApplyStmt{
// 				Kind: lexer.X,
// 				QArgs: []ast.IdentExpr{
// 					{
// 						Value: "q",
// 					},
// 				},
// 			},
// 			&ast.AssignStmt{
// 				Left: &ast.IdentExpr{
// 					Value: "c",
// 				},
// 				Right: &ast.MeasureStmt{
// 					Kind: lexer.MEASURE,
// 					QArgs: []ast.IdentExpr{
// 						{
// 							Value: "q",
// 						},
// 					},
// 				},
// 			},
// 			&ast.ArrowStmt{
// 				Left: &ast.MeasureStmt{
// 					Kind: lexer.MEASURE,
// 					QArgs: []ast.IdentExpr{
// 						{
// 							Value: "q",
// 						},
// 					},
// 				},
// 				Right: &ast.IdentExpr{
// 					Value: "c",
// 				},
// 			},
// 		},
// 	}

// 	ast.Print(p)

// 	// Output:
// 	// *ast.OpenQASM {
// 	// .  Version: 3.0
// 	// .  Incls: []string (len = 1) {
// 	// .  .  0: "stdgates.qasm"
// 	// .  }
// 	// .  Stmts: []ast.Stmt (len = 5) {
// 	// .  .  0: *ast.DeclStmt {
// 	// .  .  .  Decl: *ast.GenDecl {
// 	// .  .  .  .  Kind: qubit
// 	// .  .  .  .  Name: *ast.IdentExpr {
// 	// .  .  .  .  .  Value: q
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  }
// 	// .  .  1: *ast.ExprStmt {
// 	// .  .  .  X: *ast.ResetExpr {
// 	// .  .  .  .  QArgs: ast.IdentList {
// 	// .  .  .  .  .  List: []ast.IdentExpr (len = 1) {
// 	// .  .  .  .  .  .  0: ast.IdentExpr {
// 	// .  .  .  .  .  .  .  Value: q
// 	// .  .  .  .  .  .  }
// 	// .  .  .  .  .  }
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  }
// 	// .  .  2: *ast.ApplyStmt {
// 	// .  .  .  Kind: x
// 	// .  .  .  QArgs: []ast.IdentExpr (len = 1) {
// 	// .  .  .  .  0: ast.IdentExpr {
// 	// .  .  .  .  .  Value: q
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  }
// 	// .  .  3: *ast.AssignStmt {
// 	// .  .  .  Left: *ast.IdentExpr {
// 	// .  .  .  .  Value: c
// 	// .  .  .  }
// 	// .  .  .  Right: *ast.MeasureStmt {
// 	// .  .  .  .  Kind: measure
// 	// .  .  .  .  QArgs: []ast.IdentExpr (len = 1) {
// 	// .  .  .  .  .  0: ast.IdentExpr {
// 	// .  .  .  .  .  .  Value: q
// 	// .  .  .  .  .  }
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  }
// 	// .  .  4: *ast.ArrowStmt {
// 	// .  .  .  Left: *ast.MeasureStmt {
// 	// .  .  .  .  Kind: measure
// 	// .  .  .  .  QArgs: []ast.IdentExpr (len = 1) {
// 	// .  .  .  .  .  0: ast.IdentExpr {
// 	// .  .  .  .  .  .  Value: q
// 	// .  .  .  .  .  }
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  .  Right: *ast.IdentExpr {
// 	// .  .  .  .  Value: c
// 	// .  .  .  }
// 	// .  .  }
// 	// .  }
// 	// }
// }

// func ExamplePrint_gate() {
// 	g := &ast.GateStmt{
// 		Kind: lexer.GATE,
// 		Name: "bell",
// 		QArgs: []ast.IdentExpr{
// 			{
// 				Value: "q0",
// 			},
// 			{
// 				Value: "q1",
// 			},
// 		},
// 		Statements: []ast.Stmt{
// 			&ast.ApplyStmt{
// 				Kind: lexer.H,
// 				QArgs: []ast.IdentExpr{
// 					{
// 						Value: "q0",
// 					},
// 				},
// 			},
// 			&ast.ApplyStmt{
// 				Kind: lexer.CX,
// 				QArgs: []ast.IdentExpr{
// 					{
// 						Value: "q0",
// 					},
// 					{
// 						Value: "q1",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	ast.Print(g)

// 	// Output:
// 	// *ast.GateStmt {
// 	// .  Kind: gate
// 	// .  Name: bell
// 	// .  QArgs: []ast.IdentExpr (len = 2) {
// 	// .  .  0: ast.IdentExpr {
// 	// .  .  .  Value: q0
// 	// .  .  }
// 	// .  .  1: ast.IdentExpr {
// 	// .  .  .  Value: q1
// 	// .  .  }
// 	// .  }
// 	// .  Statements: []ast.Stmt (len = 2) {
// 	// .  .  0: *ast.ApplyStmt {
// 	// .  .  .  Kind: h
// 	// .  .  .  QArgs: []ast.IdentExpr (len = 1) {
// 	// .  .  .  .  0: ast.IdentExpr {
// 	// .  .  .  .  .  Value: q0
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  }
// 	// .  .  1: *ast.ApplyStmt {
// 	// .  .  .  Kind: cx
// 	// .  .  .  QArgs: []ast.IdentExpr (len = 2) {
// 	// .  .  .  .  0: ast.IdentExpr {
// 	// .  .  .  .  .  Value: q0
// 	// .  .  .  .  }
// 	// .  .  .  .  1: ast.IdentExpr {
// 	// .  .  .  .  .  Value: q1
// 	// .  .  .  .  }
// 	// .  .  .  }
// 	// .  .  }
// 	// .  }
// 	// }
// }
