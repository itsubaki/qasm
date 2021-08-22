# qasm

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/qasm)](https://pkg.go.dev/github.com/itsubaki/qasm)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/qasm?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/qasm)
[![tests](https://github.com/itsubaki/qasm/workflows/tests/badge.svg?branch=main)](https://github.com/itsubaki/qasm/actions)

 - Run Quantum Computation Simulator with OpenQASM 3.0

## Install

```shell
go install github.com/itsubaki/qasm@latest
```

## Example

```shell
$ qasm -f testdata/bell.qasm
OPENQASM 3.0;
include "itsubaki/q.qasm";

qubit[2] q;
reset q;

h q[0];
cx q[0], q[1];

[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
```

## REPL

```shell
$ qasm repl
>> OPENQASM 3.0;
>> include "itsubaki/q.qasm";
>> 
>> qubit q;
[0][  0]( 1.0000 0.0000i): 1.0000
>> qubit p;
[0 0][  0   0]( 1.0000 0.0000i): 1.0000
>> h q;
[0 0][  0   0]( 0.7071 0.0000i): 0.5000
[1 0][  1   0]( 0.7071 0.0000i): 0.5000
>> cx q, p;
[0 0][  0   0]( 0.7071 0.0000i): 0.5000
[1 1][  1   1]( 0.7071 0.0000i): 0.5000
>> quit
```

## AST

```shell
$ qasm parse -f testdata/bell_gate.qasm 
*ast.OpenQASM {
.  Version: 3.0
.  Incls: []ast.Stmt (len = 1) {
.  .  0: *ast.InclStmt {
.  .  .  Path: *ast.IdentExpr {
.  .  .  .  Value: "itsubaki/q.qasm"
.  .  .  }
.  .  }
.  }
.  Stmts: []ast.Stmt (len = 4) {
.  .  0: *ast.DeclStmt {
.  .  .  Decl: *ast.GateDecl {
.  .  .  .  Name: bell
.  .  .  .  Params: ast.ParenExpr {
.  .  .  .  .  List: ast.ExprList {}
.  .  .  .  }
.  .  .  .  QArgs: ast.ExprList {
.  .  .  .  .  List: []ast.Expr (len = 2) {
.  .  .  .  .  .  0: *ast.IdentExpr {
.  .  .  .  .  .  .  Value: q0
.  .  .  .  .  .  }
.  .  .  .  .  .  1: *ast.IdentExpr {
.  .  .  .  .  .  .  Value: q1
.  .  .  .  .  .  }
.  .  .  .  .  }
.  .  .  .  }
.  .  .  .  Body: *ast.BlockStmt {
.  .  .  .  .  List: []ast.Stmt (len = 2) {
.  .  .  .  .  .  0: *ast.ExprStmt {
.  .  .  .  .  .  .  X: *ast.ApplyExpr {
.  .  .  .  .  .  .  .  Kind: h
.  .  .  .  .  .  .  .  Params: ast.ParenExpr {
.  .  .  .  .  .  .  .  .  List: ast.ExprList {}
.  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  QArgs: ast.ExprList {
.  .  .  .  .  .  .  .  .  List: []ast.Expr (len = 1) {
.  .  .  .  .  .  .  .  .  .  0: *ast.IdentExpr {
.  .  .  .  .  .  .  .  .  .  .  Value: q0
.  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  }
.  .  .  .  .  .  }
.  .  .  .  .  .  1: *ast.ExprStmt {
.  .  .  .  .  .  .  X: *ast.ApplyExpr {
.  .  .  .  .  .  .  .  Kind: cx
.  .  .  .  .  .  .  .  Params: ast.ParenExpr {
.  .  .  .  .  .  .  .  .  List: ast.ExprList {}
.  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  QArgs: ast.ExprList {
.  .  .  .  .  .  .  .  .  List: []ast.Expr (len = 2) {
.  .  .  .  .  .  .  .  .  .  0: *ast.IdentExpr {
.  .  .  .  .  .  .  .  .  .  .  Value: q0
.  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  .  1: *ast.IdentExpr {
.  .  .  .  .  .  .  .  .  .  .  Value: q1
.  .  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  .  }
.  .  .  .  .  .  .  }
.  .  .  .  .  .  }
.  .  .  .  .  }
.  .  .  .  }
.  .  .  }
.  .  }
.  .  1: *ast.DeclStmt {
.  .  .  Decl: *ast.GenDecl {
.  .  .  .  Kind: qubit
.  .  .  .  Type: *ast.IndexExpr {
.  .  .  .  .  Name: *ast.IdentExpr {
.  .  .  .  .  .  Value: qubit
.  .  .  .  .  }
.  .  .  .  .  Value: 2
.  .  .  .  }
.  .  .  .  Name: *ast.IdentExpr {
.  .  .  .  .  Value: q
.  .  .  .  }
.  .  .  }
.  .  }
.  .  2: *ast.ExprStmt {
.  .  .  X: *ast.ResetExpr {
.  .  .  .  QArgs: ast.ExprList {
.  .  .  .  .  List: []ast.Expr (len = 1) {
.  .  .  .  .  .  0: *ast.IdentExpr {
.  .  .  .  .  .  .  Value: q
.  .  .  .  .  .  }
.  .  .  .  .  }
.  .  .  .  }
.  .  .  }
.  .  }
.  .  3: *ast.ExprStmt {
.  .  .  X: *ast.CallExpr {
.  .  .  .  Name: bell
.  .  .  .  Params: ast.ParenExpr {
.  .  .  .  .  List: ast.ExprList {}
.  .  .  .  }
.  .  .  .  QArgs: ast.ExprList {
.  .  .  .  .  List: []ast.Expr (len = 2) {
.  .  .  .  .  .  0: *ast.IndexExpr {
.  .  .  .  .  .  .  Name: *ast.IdentExpr {
.  .  .  .  .  .  .  .  Value: q
.  .  .  .  .  .  .  }
.  .  .  .  .  .  .  Value: 0
.  .  .  .  .  .  }
.  .  .  .  .  .  1: *ast.IndexExpr {
.  .  .  .  .  .  .  Name: *ast.IdentExpr {
.  .  .  .  .  .  .  .  Value: q
.  .  .  .  .  .  .  }
.  .  .  .  .  .  .  Value: 1
.  .  .  .  .  .  }
.  .  .  .  .  }
.  .  .  .  }
.  .  .  }
.  .  }
.  }
}
```

## built-in

 * `x`, `y`, `z`, `h`, `s`, `t`
 * `cx`, `cz`
 * `ccx`
 * `swap`, `qft`, `iqft`
 * `cmodexp2`
 * `print`
