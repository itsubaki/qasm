OPENQASM 3.0;
include "testdata/stdgates.qasm";

gate bell q, p { h q; cx q, p; }

qubit[2] q;
reset q;

bell q[0], q[1];
