# qasm
Run Quantum Computation Simulator with OpenQASM 3.0


## Install

```shell
go install github.com/itsubaki/qasm@latest
```

## Example

```shell
$ qasm -f testdata/bell.qasm
OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
reset q;

h q[0];
cx q[0], q[1];

[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
```

## REPL

```shell
$ qasm
>> qubit q
>> h q
>> print
[0][  0]( 0.7071 0.0000i): 0.5000
[1][  1]( 0.7071 0.0000i): 0.5000
```

 * `print` is `built-in` operation
