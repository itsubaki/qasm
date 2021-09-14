OPENQASM 3.0;
include "testdata/gate.qasm";

qubit[2] q;
bit[2]   c;
reset q;

H  q[0];
CX q[0], q[1];
