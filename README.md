# qasm
Run Quantum Computation Simulator with OpenQASM 3.0


## Install

```shell
go install github.com/itsubaki/qasm@latest
```

## Example

```shell
$ cat testdata/bell.qasm 
OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2]   c;

reset q;

h  q[0];
cx q[0], q[1];

c[0] = measure q[0];
c[1] = measure q[1];
```

```shell
$ qasm lex -f testdata/bell.qasm 
OPENQASM FLOAT ; 
INCLUDE STRING ; 
QUBIT [ INT ] IDENT ; 
BIT [ INT ] IDENT ; 
RESET IDENT ; 
H IDENT [ INT ] ; 
CX IDENT [ INT ] , IDENT [ INT ] ; 
IDENT [ INT ] = MEASURE IDENT [ INT ] ; 
IDENT [ INT ] = MEASURE IDENT [ INT ] ; 
```

