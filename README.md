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
>> gate X q { U(pi, 0.0, pi) q; }
>> gate CX q, p { ctrl @ X q, p; }                             
>> 
>> qubit[2] q;
[00][  0]( 1.0000 0.0000i): 1.0000
>> h q[0];
[00][  0]( 0.7071 0.0000i): 0.5000
[10][  2]( 0.7071 0.0000i): 0.5000
>> CX q[0], q[1];
[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
```

## built-in

 * `U`, `x`, `y`, `z`, `h`, `s`, `t`
 * `cx`, `cz`
 * `ccx`
 * `swap`, `qft`, `iqft`
 * `cmodexp2`
 * `print`
