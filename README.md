# qasm

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/qasm)](https://pkg.go.dev/github.com/itsubaki/qasm)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/qasm?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/qasm)
[![tests](https://github.com/itsubaki/qasm/workflows/tests/badge.svg)](https://github.com/itsubaki/qasm/actions)
[![codecov](https://codecov.io/gh/itsubaki/qasm/branch/main/graph/badge.svg?token=94KAQTK9KT)](https://codecov.io/gh/itsubaki/qasm)

Quantum computing simulator for OpenQASM.

## Installation

```shell
go install github.com/itsubaki/qasm@latest
```

```shell
% qasm -help
Usage of qasm:
  -f string
        filepath
  -lex
        Lex the input into a sequence of tokens
  -parse
        Parse the input and convert it into an AST (abstract syntax tree)
  -repl
        REPL(read-eval-print loop) mode
  -verbose
        Enable verbose output
```

## Examples

```shell
% qasm < testdata/bell.qasm
[00] ( 0.7071 0.0000i): 0.5000
[11] ( 0.7071 0.0000i): 0.5000
```

```shell
% qasm -repl
qasm> OPENQASM 3.0;
qasm> 
qasm> qubit[2] q;
qasm> U(pi/2, 0, pi) q[0];
qasm> ctrl @ U(pi, 0, pi) q[0], q[1];
qasm> :print
--- STATE ---
[00] ( 0.7071 0.0000i): 0.5000
[11] ( 0.7071 0.0000i): 0.5000
--- ENVIRONMENT ---
const     : map[]
variable  : map[]
bit       : map[]
bit[]     : map[]
qubit     : map[q:[0 1]]
gate      : []
subroutine: []
```

```shell
% qasm -repl
qasm> OPENQASM 3.0;
qasm> 
qasm> const float ratio = pi;
qasm> int n = 2;
qasm> if (n > 0) { n = n*ratio; }
qasm> :print
--- STATE ---
--- ENVIRONMENT ---
const     : map[ratio:3.141592653589793]
variable  : map[n:6.283185307179586]
bit       : map[]
qubit     : map[]
gate      : []
subroutine: []
```
