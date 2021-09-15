OPENQASM 3.0;
include "testdata/gate.qasm";

qubit[2] q;
reset q;

h  q[0];
cx q[0], q[1];
