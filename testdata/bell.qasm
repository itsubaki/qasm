OPENQASM 3.0;
include "itsubaki/q.qasm";

qubit q[2];
reset q;

h  q[0];
cx q[0], q[1];
