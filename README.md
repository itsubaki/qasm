# qasm
Run Quantum Computation Simulator with OpenQASM 3.0


## Install

```shell
go install github.com/itsubaki/qasm@latest
```

## Example

```shell
$ qasm
>> qubit q
>> h q
>> print
[0][  0]( 0.7071 0.0000i): 0.5000
[1][  1]( 0.7071 0.0000i): 0.5000
```

```shell
$ qasm -f testdata/print.qasm
OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
reset q;

h q;
print;

[0][  0]( 0.7071 0.0000i): 0.5000
[1][  1]( 0.7071 0.0000i): 0.5000
```
