OPENQASM 3.0;
include "testdata/stdgates.qasm";

qubit[2] q;

h q[0];
cx q[0], q[1];
