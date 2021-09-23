OPENQASM 3.0;
include "testdata/stdgates.qasm";

gate bell q, p { h q; cx q, p; }

qubit c;
qubit[2] q;
reset q;

x c;
ctrl @ bell c, q[0], q[1];
