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
$ qasm
>> OPENQASM 3.0;
>> include "itsubaki/q.qasm";
>> 
>> qubit q
[0][  0]( 1.0000 0.0000i): 1.0000
>> qubit p
[00][  0]( 1.0000 0.0000i): 1.0000
>> h q
[00][  0]( 0.7071 0.0000i): 0.5000
[10][  2]( 0.7071 0.0000i): 0.5000
>> cx q, p
[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
>> quit
```

## built-in

 * `print`
 * `x`, `y`, `z`, `h`, `s`, `t`
 * `cx`, `cz`
 * `ccx`
 * `swap`, `qft`, `iqft`
