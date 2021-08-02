# qasm
Run Quantum Computation Simulator with OpenQASM 3.0


## Install

```shell
go install github.com/itsubaki/qasm@latest
```

## Example

```shell
$ qasm lex -f pkg/_testdata/bell.qasm 
OPENQASM OPENQASM 3.0 FLOAT ; ; 
include INCLUDE "stdgates.qasm" STRING ; ; 
qubit QUBIT [ [ 2 INT ] ] q IDENT ; ; 
bit BIT [ [ 2 INT ] ] c IDENT ; ; 
reset RESET q IDENT ; ; 
h H q IDENT [ [ 0 INT ] ] ; ; 
cx CX q IDENT [ [ 0 INT ] ] , , q IDENT [ [ 1 INT ] ] ; ; 
c IDENT [ [ 0 INT ] ] = = measure MEASURE q IDENT [ [ 0 INT ] ] ; ; 
c IDENT [ [ 1 INT ] ] = = measure MEASURE q IDENT [ [ 1 INT ] ] ; ; 
```
