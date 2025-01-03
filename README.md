# qasm

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/qasm)](https://pkg.go.dev/github.com/itsubaki/qasm)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/qasm?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/qasm)
[![tests](https://github.com/itsubaki/qasm/workflows/tests/badge.svg)](https://github.com/itsubaki/qasm/actions)
[![codecov](https://codecov.io/gh/itsubaki/qasm/branch/main/graph/badge.svg?token=94KAQTK9KT)](https://codecov.io/gh/itsubaki/qasm)

- Quantum Computation Simulator with [OpenQASM 3.0](https://openqasm.com)
- Currently rewriting using ANTLR4 [WIP](https://github.com/itsubaki/qasm/tree/antlr)

## Install

```shell
go install github.com/itsubaki/qasm@latest
```

## Example

```shell
$ cat testdata/bell.qasm
OPENQASM 3.0;
include "testdata/stdgates.qasm";

qubit[2] q;
reset q;

h q[0];
cx q[0], q[1];
```

```shell
$ qasm -f testdata/bell.qasm
[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
```

## REPL

```shell
$ qasm repl
>> OPENQASM 3.0;
>>
>> gate h q { U(pi/2.0, 0.0, pi) q; }
>> gate x q { U(pi, 0, pi) q; }
>> gate cx q, p { ctrl @ x q, p; }
>>
>> qubit[2] q;
[00][  0]( 1.0000 0.0000i): 1.0000
>> h q[0];
[00][  0]( 0.7071 0.0000i): 0.5000
[10][  2]( 0.7071 0.0000i): 0.5000
>> cx q[0], q[1];
[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
```

## built-in

- `U`, `X`, `Y`, `Z`, `H`, `S`, `T`
- `QFT`, `IQFT`, `CMODEXP2`
- `print`
