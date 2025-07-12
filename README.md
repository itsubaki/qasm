# qasm

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/qasm)](https://pkg.go.dev/github.com/itsubaki/qasm)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/qasm?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/qasm)
[![tests](https://github.com/itsubaki/qasm/workflows/tests/badge.svg)](https://github.com/itsubaki/qasm/actions)
[![codecov](https://codecov.io/gh/itsubaki/qasm/branch/main/graph/badge.svg?token=94KAQTK9KT)](https://codecov.io/gh/itsubaki/qasm)

- Quantum Computation Simulator with [OpenQASM 3.0](https://openqasm.com)

## Installation

```shell
go install github.com/itsubaki/qasm@latest
```

## Examples

```shell
% qasm < testdata/bell.qasm
[00][  0]( 0.7071 0.0000i): 0.5000
[11][  3]( 0.7071 0.0000i): 0.5000
const     : map[]
variable  : map[]
bit       : map[]
qubit     : map[q:[0 1]]
gate      : [cx h i x y z]
subroutine: []
```

```shell
% go run cmd/repl/main.go
>> OPENQASM 3.0;
>> 
>> qubit[2] q;
[00][  0]( 1.0000 0.0000i): 1.0000
>> U(pi/2, 0, pi) q[0];
[00][  0]( 0.7071 0.0000i): 0.5000
[10][  2]( 0.7071 0.0000i): 0.5000
>> ctrl @ U(pi, 0, pi) q[0], q[1];
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
>> print;
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
[@-1,15:19='qubit',<31>,3:0]
[@-1,20:20='[',<56>,3:5]
[@-1,21:21='2',<91>,3:6]
[@-1,22:22=']',<57>,3:7]
[@-1,24:24='q',<93>,3:9]
[@-1,25:25=';',<63>,3:10]
[@-1,27:27='U',<93>,4:0]
[@-1,28:28='(',<60>,4:1]
[@-1,29:30='pi',<93>,4:2]
[@-1,31:31='/',<73>,4:4]
[@-1,32:34='2.0',<95>,4:5]
[@-1,35:35=',',<65>,4:8]
[@-1,37:37='0',<91>,4:10]
[@-1,38:38=',',<65>,4:11]
[@-1,40:41='pi',<93>,4:13]
[@-1,42:42=')',<61>,4:15]
[@-1,47:47='q',<93>,4:20]
[@-1,48:48='[',<56>,4:21]
[@-1,49:49='0',<91>,4:22]
[@-1,50:50=']',<57>,4:23]
[@-1,51:51=';',<63>,4:24]
[@-1,53:56='ctrl',<47>,5:0]
[@-1,58:58='@',<80>,5:5]
[@-1,60:60='U',<93>,5:7]
[@-1,61:61='(',<60>,5:8]
[@-1,62:63='pi',<93>,5:9]
[@-1,64:64=',',<65>,5:11]
[@-1,66:66='0',<91>,5:13]
[@-1,67:67=',',<65>,5:14]
[@-1,69:70='pi',<93>,5:16]
[@-1,71:71=')',<61>,5:18]
[@-1,73:73='q',<93>,5:20]
[@-1,74:74='[',<56>,5:21]
[@-1,75:75='0',<91>,5:22]
[@-1,76:76=']',<57>,5:23]
[@-1,77:77=',',<65>,5:24]
[@-1,79:79='q',<93>,5:26]
[@-1,80:80='[',<56>,5:27]
[@-1,81:81='1',<91>,5:28]
[@-1,82:82=']',<57>,5:29]
[@-1,83:83=';',<63>,5:30]
```

```shell
% go run cmd/parse/main.go < testdata/bell.qasm
(program
  (version OPENQASM 3.0 ;)
  (statementOrScope
    (statement
      (quantumDeclarationStatement
        (qubitType qubit
          (designator [ (expression 2) ])
        ) q;
      )
    )
  )
  (statementOrScope
    (statement
      (gateCallStatement U
        (
          (expressionList
            (expression (expression pi) / (expression 2.0)),
            (expression 0),
            (expression pi)
          )
        )
        (gateOperandList
          (gateOperand
            (indexedIdentifier q (indexOperator [ (expression 0) ]))
          )
        );
      )
    )
  )
  (statementOrScope
    (statement
      (gateCallStatement
        (gateModifier ctrl @) U
        (
          (expressionList
            (expression pi),
            (expression 0),
            (expression pi)
          )
        )
        (gateOperandList
          (gateOperand
            (indexedIdentifier q (indexOperator [ (expression 0) ]))
          ),
          (gateOperand
            (indexedIdentifier q (indexOperator [ (expression 1) ]))
          )
        );
      )
    )
  )
  <EOF>
)
```
