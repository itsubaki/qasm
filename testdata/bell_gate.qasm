OPENQASM 3.0;
include "testdata/gate.qasm";

qubit[2] q;
bit[2]   c;
reset q;

bell q[0], q[1];
measure q -> c;
