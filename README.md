# qasm

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/qasm)](https://pkg.go.dev/github.com/itsubaki/qasm)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/qasm?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/qasm)
[![tests](https://github.com/itsubaki/qasm/workflows/tests/badge.svg)](https://github.com/itsubaki/qasm/actions)
[![codecov](https://codecov.io/gh/itsubaki/qasm/branch/main/graph/badge.svg?token=94KAQTK9KT)](https://codecov.io/gh/itsubaki/qasm)

- Quantum Computation Simulator with [OpenQASM 3.0](https://openqasm.com)

## Example

```shell
% go run cmd/repl/main.go                       
>> OPENQASM 3.0;
>> include "testdata/stdgates.qasm";
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

```shell
% go run cmd/repl/main.go
>> OPENQASM 3.0;
>> 
>> const float ratio = pi;
>> 
>> int n = 2;
>> if (n > 0) { n = n*ratio; }
>>
>> print
const     : map[ratio:3.141592653589793]
variable  : map[n:6.283185307179586]
bit       : map[]
qubit     : map[]
gate      : []
subroutine: []
```

```shell
% go run cmd/lex/main.go < testdata/bell.qasm
[@-1,0:7='OPENQASM',<1>,1:0]
[@-1,9:11='3.0',<103>,1:9]
[@-1,12:12=';',<63>,1:12]
[@-1,14:20='include',<2>,2:0]
[@-1,22:46='"testdata/stdgates.qasm"',<105>,2:8]
[@-1,47:47=';',<63>,2:33]
[@-1,50:54='qubit',<31>,4:0]
[@-1,55:55='[',<56>,4:5]
[@-1,56:56='2',<91>,4:6]
[@-1,57:57=']',<57>,4:7]
[@-1,59:59='q',<93>,4:9]
[@-1,60:60=';',<63>,4:10]
[@-1,63:63='h',<93>,6:0]
[@-1,65:65='q',<93>,6:2]
[@-1,66:66='[',<56>,6:3]
[@-1,67:67='0',<91>,6:4]
[@-1,68:68=']',<57>,6:5]
[@-1,69:69=';',<63>,6:6]
[@-1,71:72='cx',<93>,7:0]
[@-1,74:74='q',<93>,7:3]
[@-1,75:75='[',<56>,7:4]
[@-1,76:76='0',<91>,7:5]
[@-1,77:77=']',<57>,7:6]
[@-1,78:78=',',<65>,7:7]
[@-1,80:80='q',<93>,7:9]
[@-1,81:81='[',<56>,7:10]
[@-1,82:82='1',<91>,7:11]
[@-1,83:83=']',<57>,7:12]
[@-1,84:84=';',<63>,7:13]
```

```shell
go run cmd/parse/main.go < testdata/bell.qasm
(program
  (version OPENQASM 3.0 ;)
  (statementOrScope
    (statement
      (includeStatement include "testdata/stdgates.qasm" ;)))
  (statementOrScope
    (statement
      (quantumDeclarationStatement
        (qubitType qubit (designator [ (expression 2) ])) q ;)))
  (statementOrScope
    (statement
      (gateCallStatement h
        (gateOperandList
          (gateOperand
            (indexedIdentifier q (indexOperator [ (expression 0) ])))
        ;))))
  (statementOrScope
    (statement
      (gateCallStatement cx
        (gateOperandList
          (gateOperand
            (indexedIdentifier q (indexOperator [ (expression 0) ]))),
          (gateOperand
            (indexedIdentifier q (indexOperator [ (expression 1) ])))
        ;))))
  <EOF>)
```
